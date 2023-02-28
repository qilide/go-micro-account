# Account Service

This is the Account service

Generated with

```
micro new --namespace=go.micro --type=service account
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.service.account
- Type: service
- Alias: account

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./account-service
```

Build a docker image
```
make docker
```

### 启动命令
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