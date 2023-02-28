package redis

import (
	"account/common/micro"
	"fmt"
	"github.com/go-redis/redis"
)

// Rdb 声明一个全局的rdb变量
var Rdb *redis.Client

// Init 初始化连接
func Init() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			micro.ConsulInfo.Redis.Host,
			micro.ConsulInfo.Redis.Port,
		),
		DB:       int(micro.ConsulInfo.Redis.Db), // use default DB
		PoolSize: int(micro.ConsulInfo.Redis.PoolSize),
	})

	_, err = Rdb.Ping().Result()
	return err
}

func Close() {
	_ = Rdb.Close()
}
