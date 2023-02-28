package main

import (
	micro2 "account/common/micro"
	"account/config/logger"
	"account/config/mysql"
	"account/config/redis"
	"account/domain/repository"
	"account/domain/service"
	"account/handler"
	"account/proto/account"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
)

func main() {
	// 1.配置中心
	consulConfig, err := micro2.GetConsulConfig("localhost", 8500, "/micro/config")
	if err != nil {
		fmt.Printf("Init consulConfig failed, err: %v\n", err)
	}
	// 2.注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	if err := micro2.GetAccountFromConsul(consulConfig, "account"); err != nil {
		fmt.Printf("Init consul failed, err: %v\n", err)
	}
	fmt.Println(micro2.ConsulInfo)
	// 3.zap日志初始化
	if err := logger.Init(); err != nil {
		fmt.Printf("Init logger failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync()
	// 4.jaeger 链路追踪
	t, io, err := micro2.NewTracer(micro2.ConsulInfo.Jaeger.ServiceName, micro2.ConsulInfo.Jaeger.Addr)
	if err != nil {
		logger.Error(err)
		return
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	// 5.初始化数据库
	db, err := mysql.MysqlInit(micro2.ConsulInfo.Mysql.User, micro2.ConsulInfo.Mysql.Pwd, micro2.ConsulInfo.Mysql.Database)
	if err != nil {
		logger.Error(err)
		return
	}
	defer db.Close()
	// 创建实例
	accountService := service.NewUserService(repository.NewUserRepository(db))
	// 6.初始化Redis连接
	if err := redis.Init(); err != nil {
		logger.Error(err)
		return
	}
	defer redis.Close()
	// 7.暴露监控地址
	micro2.PrometheusBoot(micro2.ConsulInfo.Prometheus.Host, int(micro2.ConsulInfo.Prometheus.Port))
	// 8.注册服务
	registryService := micro.NewService(
		micro.Name(micro2.ConsulInfo.Micro.Name),
		micro.Version(micro2.ConsulInfo.Micro.Version),
		//暴露的服务地址
		micro.Address(micro2.ConsulInfo.Micro.Address),
		//添加consul 注册中心
		micro.Registry(consulRegistry),
		//添加链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(int(micro2.ConsulInfo.Ratelimit.QPS))),
		//添加监控
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)
	// 9.初始化服务
	registryService.Init()
	// 10.注册Handle
	account.RegisterAccountHandler(registryService.Server(), &handler.Account{AccountService: accountService})
	// 11.启动服务
	if err := registryService.Run(); err != nil {
		logger.Fatal(err)
	}
}
