package micro

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-plugins/config/source/consul/v2"
	"strconv"
)

// GetConsulConfig 设置配置中心
func GetConsulConfig(host string, port int64, prefix string) (config.Config, error) {
	consulSource := consul.NewSource(
		//设置配置中心的地址
		consul.WithAddress(host+":"+strconv.FormatInt(port, 10)),
		//设置前缀，不设置默认前缀 /micro/config
		consul.WithPrefix(prefix),
		//是否移除前缀，这里是设置为true，表示可以不带前缀直接获取对应配置
		consul.StripPrefix(true),
	)
	//配置初始化
	newConfig, err := config.NewConfig()
	if err != nil {
		return newConfig, err
	}
	//加载配置
	err = newConfig.Load(consulSource)
	return newConfig, err
}

type Account struct {
	Name    string `json:"name"`
	Title   string `json:"title"`
	Mode    string `json:"mode"`
	Port    int64  `json:"port"`
	Version string `json:"version"`
}

type Mysql struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Pwd      string `json:"pwd"`
	Database string `json:"database"`
	Port     int64  `json:"port"`
}

type Log struct {
	Level string `json:"level"`
	Filename string `json:"filename"`
	MaxSize int64 `json:"max_size"`
	MaxAge int64 `json:"max_age"`
	MaxBackips int64 `json:"max_backips"`
}

type Redis struct {
	Host string `json:"host"`
	Port int64 `json:"port"`
	Password string `json:"password"`
	Db int64 `json:"db"`
	PoolSize int64 `json:"pool_size"`
}

type Email struct {
	User string `json:"user"`
	Pass string `json:"pass"`
	Host string `json:"host"`
	Port int64 `json:"port"`
	Rename string `json:"rename"`
}

type Consul struct {
	Host string `json:"host"`
	Port int64 `json:"port"`
	Prefix string `json:"prefix"`
	ConsulRegistry string `json:"consulRegistry"`
}

type Jaeger struct {
	ServiceName string `json:"serviceName"`
	Addr string `json:"addr"`
}

type Prometheus struct {
	Host string `json:"host"`
	Port int64 `json:"port"`
}

type Ratelimit struct {
	QPS int64 `json:"QPS"`
}

type Micro struct {
	Name string `json:"name"`
	Version string `json:"version"`
	Address string `json:"address"`
}

type ConsulConfig struct {
	Account Account `json:"account"`
	Mysql Mysql `json:"mysql"`
	Log Log `json:"log"`
	Redis Redis `json:"redis"`
	Email Email `json:"email"`
	Consul Consul `json:"consul"`
	Jaeger Jaeger `json:"jaeger"`
	Prometheus Prometheus `json:"prometheus"`
	Ratelimit Ratelimit `json:"ratelimit"`
	Micro Micro `json:"micro"`
}

var(
	ConsulInfo *ConsulConfig
)

// GetAccountFromConsul 获取 consul 的配置
func GetAccountFromConsul(config config.Config, path ...string) error {
	consulData := &ConsulConfig{}
	config.Get(path...).Scan(consulData)
	ConsulInfo = consulData
	return nil
}
