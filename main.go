package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	kafkago "github.com/segmentio/kafka-go"
	"os"
	"strings"

	"github.com/xchenhao/data-transport-service/config"
	mykafka "github.com/xchenhao/data-transport-service/kafka"
	mongo_msg "github.com/xchenhao/data-transport-service/kafka/mongo-msg"
	"github.com/xchenhao/data-transport-service/logger"
	"github.com/xchenhao/data-transport-service/sql"
)

func parseArgs() string {
	configFile := flag.String("config", "", "config file path")
	showHelp := flag.Bool("help", false, "show help message")
	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *configFile == "" {
		fmt.Fprintln(os.Stderr, "Please specify config file path")
		flag.Usage()
		os.Exit(1)
	}

	return *configFile
}

var (
	mappings map[string]*MongoCollectionToSQLDBTable
	db *gorm.DB
)

func main() {
	configFilePath := parseArgs()

	conf, err := config.LoadConfig(configFilePath)
	if err != nil {
		logger.Fatalln(err)
	}

	db, err = sql.InitDB(conf.DB)
	if err != nil {
		logger.Fatalln(err)
	}

	mappings, err = LoadMapping(conf.MongoDBMapToSQLFile)
	if err != nil {
		logger.Fatalln(err)
	}

	kafkaReader := mykafka.Reader(conf.Kafka)
	defer kafkaReader.Close()

	fmt.Println("Data Transport Service Running...")
	for {
		ctx := context.Background()
		message, err := kafkaReader.FetchMessage(context.Background())
		if err != nil {
			logger.Fatalln(err)
		}

		func() {
			if err := recover(); err != nil {
				logger.Println("consume message panic: ", err)
			}

			err = consumeMessageDispatch(message, func() error {
				return kafkaReader.CommitMessages(ctx, message)
			})
			if err != nil {
				logger.Println("consume message error: " + err.Error())
			}
		}()
	}
}

func consumeMessageDispatch(message kafkago.Message, callback func() error) error {
	if message.Value == nil {  // 删除操作有两条消息，其中一条 body 为空
		return nil
	}
	msg := new(mongo_msg.MongoKafkaMsgBody)
	err := json.Unmarshal(message.Value, msg)
	if err != nil {
		return errors.New("message value decode error: " + err.Error())
	}
	if !msg.IsFromMongoDB() {
		return errors.New("only support mongodb message")
	}
	collection := msg.Payload.Source.Collection
	mappingRule := FindItemByCollection(mappings, collection)

	switch msg.Payload.Op {
	case mongo_msg.OpDelete:
		fallthrough
	case mongo_msg.OpUpdate:
		msgKey := new(mongo_msg.MongoKafkaMsgKey)
		err := json.Unmarshal(message.Key, msgKey)
		if err != nil {
			return errors.New("message key decode error: " + err.Error())
		}
		// https://www.mongodb.com/docs/manual/core/document/
		// https://debezium.io/documentation/reference/stable/connectors/mongodb.html
		toMap, err := mongo_msg.DecodeMsgBodyFieldsToMap(msgKey.Payload)
		if err != nil {
			return errors.New("decode key fields error: " + err.Error())
		}
		if msg.Payload.Op == mongo_msg.OpDelete {
			err = db.Table(mappingRule.Table).Where(mappingRule.ColumnMapping["_id"]+" = ?", toMap["id"]).Delete(nil).Error
			if err != nil {
				return errors.New("delete record error: " + err.Error())
			}
		} else if msg.Payload.Op == mongo_msg.OpUpdate {
			upf := make(map[string]interface{})
			err := json.Unmarshal([]byte(msg.Payload.UpdateDescription.UpdatedFields), &upf)
			if err != nil {

			}
			updateFields, err := mongo_msg.DecodeMsgBodyFieldsToMap(upf)
			if err != nil {
				return errors.New("decode update fields error: " + err.Error())
			}
			updateColumns := make(map[string]interface{}, len(updateFields))
			for k, v := range updateFields {
				column, ok := mappingRule.ColumnMapping[k]
				if !ok {
					logger.Println("mongodb doc field no corresponding SQL column: " + k)
					continue
				}
				updateColumns[column] = v
			}

			err = db.Table(mappingRule.Table).
				Where(mappingRule.ColumnMapping["_id"]+" = ?", toMap["id"]).
				UpdateColumns(updateColumns).Error
			if err != nil {
				return errors.New("update record error: " + err.Error())
			}
		}
	case mongo_msg.OpInsert:
		fallthrough
	case mongo_msg.OpRead:
		after := make(map[string]interface{})
		err := json.Unmarshal([]byte(msg.Payload.After), &after)
		if err != nil {
			return errors.New("decode after fields error: " + err.Error())
		}
		mongoFieldValueMap, err := mongo_msg.DecodeMsgBodyFieldsToMap(after)
		if err != nil {
			return errors.New("decode body fields error: " + err.Error())
		}
		// columns := make([]string, 0, len(mongoFieldValueMap))
		_columns := make([]interface{}, 0, len(mongoFieldValueMap))
		values := make([]interface{}, 0, len(mongoFieldValueMap))
		for field, value := range mongoFieldValueMap {
			sqlColumn, ok := mappingRule.ColumnMapping[field]
			if !ok {
				logger.Println("mongodb doc field no corresponding SQL column: " + field)
				continue
			}
			// columns = append(columns, sqlColumn)
			_columns = append(_columns, sqlColumn)
			values = append(values, value)
		}

		sql := fmt.Sprintf("INSERT INTO `%s`", mappingRule.Table)+
			fmt.Sprintf("("+strings.TrimRight(strings.Repeat("%s, ", len(_columns)), ", ")+")", _columns...)+
			"VALUES("+strings.TrimRight(strings.Repeat("?, ", len(values)), ", ")+")"
		err = db.Exec(sql, values...).Error
		if err != nil {
			return errors.New("update record error: " + err.Error())
		}
	}

	return callback()
}
