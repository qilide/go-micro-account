### 启动命令
首先启动Docker，镜像启动可以使用docker-compose启动或者单独启动，二者选其一就可以
#### docker-compose方式启动镜像
```shell
cd docker-compose
docker-compose -f docker-compose.yml up -d
cd ..
cd docker-elk
docker-compose -f docker-stack.yml up -d
```
#### 镜像单独启动
1. 启动jaeger链路追踪：
```shell
docker run -d --name jaeger -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one
```   
2. 启动prometheus监控服务
```shell
docker run -d -p 9090:9090 bitnami/prometheus
```
3. 启动elasticsearch
```shell
docker run --name elasticsearch7.9.3 -d -e ES_JAVA_OPTS="-Xmx256m -Xms256m" --net host -e "discovery.type=single-node" -p 9200:9200 -p 9300:9300 elasticsearch:7.9.3
```
4.启动logstash
```shell
docker run -d --name logstash -p 5044:5044 -p 5000:5000 -p 9600:9600 logstash
```
5. 启动kibana
```shell
docker run --name kibana -d -p 5601:5601 kibana
```
6. 启动consul配置注册中心
```shell
docker run -d -p 8500:8500 consul
```

### consul配置中心
docker启动后，需要在配置中心consul上填写相关配置

1.访问 http://127.0.0.1:8500

2.点击Key/Value再点击Create

3.名字路径为：micro/config/account，Value格式选json

4.输入以下内容（有些参数自行修改）：
```json
{
"account":{

"name": "Account",

"title": "账号功能",

"mode": "dev",

"port": 9580,

"version": "v0.0.1"

},
"log":{

"level": "debug",

"filename": "Account.log",

"max_size": 200,

"max_age": 30,

"max_backips": 7

},
"mysql":{

"host":"127.0.0.1",

"user":"root",

"pwd":"xxx",

"database":"micro",

"port":3306

},
"redis":{

"host": "127.0.0.1",

"port": 6379,

"password": "xxx",

"db": 4,

"pool_size": 100

},
"consul":{

"host": "localhost",

"port": 8500,

"prefix": "/micro/config",

"consulRegistry": "127.0.0.1:8500"

},
"email":{

"user": "xxx@qq.com",

"pass": "xxx",

"host": "smtp.qq.com",

"port": 465,

"rename": "Account"

},
"jaeger":{

"serviceName": "go.micro.service.account",

"addr": "localhost:6831"

},
"prometheus":{

"host": "0.0.0.0",

"port": 9089

},
"ratelimit":{

"QPS": 1000

},
"micro":{

"name": "go.micro.service.account",

"version": "latest",

"address": ":9580"

}
}
```
### 数据库配置（第一次启动需做）
1.新建数据库：micro

2.进入/config/mysql/mysql.go文件

3.打开初始化表的注释语句（之后启动需要注释此语句）

### 使用
1.启动main.go

2.启动client/account.go

3.正常运行结果为打印出新添加的信息

#### 联系我们
出现问题可加群联系解决：325280438
