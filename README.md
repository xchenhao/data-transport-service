# Data Transport Service
> 基于 Debezium <br/>
> https://github.com/debezium/debezium

## 功能
将 MongoDB 数据实时同步至 MySQL

## 前置
- Zookeeper
- Kafka
- Kafka Connect
    + debezium-connector-mongodb.jar
    + mongodb-driver.jar
- MySQL

```shell
$ cd docker/debezium/kafka-connect
$ docker build -t xchenhao/kafka-connect-mongo:v1.0 . # 构建 kafka-connect 镜像

$ cd docker/debezium/
$ vim docker-compose.yml # 更改相关配置
$ vim .env # 更改相关配置
$ docker-compose up # 启动 zookeeper、kafka、kafka-connect
```

## 部署
```shell
$ go mod download
$ go build .
./data-transport-service
Please specify config file path
Usage of ./data-transport-service:
  -config string
    	config file path
  -help
    	show help message

$ ./data-transport-service -config ./config.example.yml
```
