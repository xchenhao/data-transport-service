# Data Transport Service
> 基于 Debezium <br/>
> https://github.com/debezium/debezium

## 功能
将 MongoDB 数据实时同步至 MySQL

## 项目结构概览
```shell
$ tree -L 1
.
├── config # 配置加载
├── kafka # Kafka 消息处理
├── logger # 日志相关
├── main.go # 入口文件
├── mapping.go # MongoDB 与 MySQL 列映射处理
├── sql # MySQL 配置

├── go.mod # go modules
├── go.sum
├── vendor

├── config.example.yml # 项目配置文件
├── mapping 列映射配置文件

├── Makefile # 用于构建镜像
├── Dockerfile

├── README.md # 说明
└── prerequisites # 前置依赖，与项目编译无关
```

## 前置
> 见 `prerequisites` 目录
- Zookeeper
- Kafka
- Kafka Connect
    + debezium-connector-mongodb.jar
    + mongodb-driver.jar
- MySQL

```shell
$ cd prerequisites/debezium/kafka-connect
$ docker build -t xchenhao/kafka-connect-mongo:v1.0 . # 构建 kafka-connect 镜像

$ cd prerequisites/debezium/
$ vim docker-compose.yml # 更改相关配置
$ vim .env # 更改相关配置
$ docker-compose up # 启动 zookeeper、kafka、kafka-connect
```

## 部署
- 宿主机
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

- Docker
```shell
$ docker run --rm -v $PWD/config.yml:/etc/config.yml -v $PWD/mapping:/etc/mapping  data-transport-service:0.1.0
(/go/src/github.com/xchenhao/data-transport-service/main.go:179)
[2023-04-20 11:23:26]  [43.34ms]  INSERT INTO `users`(user_id, first_name, last_name, mail)VALUES('1008', 'Li', 'Si', 'lisi@qq.com')
[1 rows affected or returned ]

(/go/src/github.com/xchenhao/data-transport-service/main.go:145)
[2023-04-20 12:01:01]  [10.33ms]  UPDATE `users` SET `last_name` = 'Wu'  WHERE (user_id = '1008')
[1 rows affected or returned ]
```