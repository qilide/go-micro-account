package main

import (
	micro2 "account/common/micro"
	"account/config/logger"
	"account/proto/account"
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/select/roundrobin/v2"
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
	newConsul := consul.NewRegistry(func(options *registry.Options) {
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
	logger.Debug("logger init success...")

	// 4.链路追踪
	t, io, err := micro2.NewTracer("go.micro.service.account", "localhost:6831")
	if err != nil {
		logger.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)
	// 5.监控
	micro2.PrometheusBoot("127.0.0.1", 9292)
	// 6.设置服务
	service := micro.NewService(
		micro.Name("go.micro.service.account.client"),
		micro.Version("latest"),
		//暴露的服务地址
		micro.Address("127.0.0.1:9580"),
		//添加注册中心
		micro.Registry(newConsul),
		//绑定链路追踪
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		//添加监控
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//作为服务端访问时生效
		micro.WrapClient(opentracing2.NewClientWrapper(opentracing.GlobalTracer())),
		//负载均衡
		micro.WrapClient(roundrobin.NewClientWrapper()),
	)
	// 7.创建服务
	accountService := account.NewAccountService("go.micro.service.account", service.Client())
	// 8.发送注册邮件
	registerMail := &account.SendMailRequest{Email: "1019528265@qq.com"}
	registerMailResponse, err := accountService.SendRegisterMail(context.TODO(), registerMail)
	if err != nil {
		logger.Error(err)
	}
	fmt.Println(registerMailResponse)
	// 9.实现注册功能
	accountAdd := &account.RegisterRequest{
		RegisterRequest: &account.UserInfoResponse{
			Username:  "夏沫の梦7777",
			FirstName: "qi66",
			Password:  "123456",
			Email:     "1019528265@qq.com",
			LastName:  "admin",
		},
		Code: registerMailResponse.Code,
	}
	registerResponse, err := accountService.Register(context.TODO(), accountAdd)
	if err != nil {
		logger.Error(err)
	}
	fmt.Println(registerResponse)
	// 10.查询用户功能
	getUser := &account.UserIdRequest{UserId: registerResponse.UserId}
	userInfoResponse, err := accountService.GetUserInfo(context.TODO(), getUser)
	if err != nil {
		logger.Error(err)
	}
	fmt.Println(userInfoResponse)
}
