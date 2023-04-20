package mongo_msg

import (
	"encoding/json"
)

type MongoDocOperation = string

const OpDelete MongoDocOperation = "d"
const OpUpdate MongoDocOperation = "u"
const OpInsert MongoDocOperation = "c"
const OpRead MongoDocOperation = "r" // 仅适用于快照

type MongoKafkaMsgKey struct {
	Schema json.RawMessage `json:"-"`
	Payload map[string]interface{} `json:"payload"`
}

// MongoKafkaMsgBody https://debezium.io/documentation/reference/stable/connectors/mongodb.html
type MongoKafkaMsgBody struct {
	Schema json.RawMessage `json:"-"`
	Payload struct{
		Before            string                                `json:"before"`
		After             string                `json:"after"`
		Source            MongoKafkaMsgBodySource        `json:"source"`
		UpdateDescription *MongoKafkaMsgBodyUpdateFields `json:"updateDescription"`
		Op                MongoDocOperation              `json:"op"`
	} `json:"payload"`
}

func (m *MongoKafkaMsgBody) IsFromMongoDB() bool {
	return m.Payload.Source.Connector == "mongodb"
}

type MongoKafkaMsgBodyUpdateFields struct {
	UpdatedFields string `json:"updatedFields"`
	RemovedFields []string `json:"removedFields"`  // todo
	TruncateArrays []string `json:"truncateArrays"`  // todo
}

type MongoKafkaMsgBodySource struct {
	Version string `json:"version"`  // connector version, eg. 2.1.4.Final
	Connector string `json:"connector"` // mongodb
	Snapshot string `json:"snapshot"` // true/false
	RS string `json:"rs"`
	DB string `json:"db"`
	Collection string `json:"collection"`
}

// DecodeMsgBodyFieldsToMap https://www.mongodb.com/docs/manual/reference/bson-types/
// https://www.mongodb.com/docs/manual/reference/mongodb-extended-json/
// TODO 改成结构体解析（需要先生成结构体代码）
func DecodeMsgBodyFieldsToMap(fieldValueMap map[string]interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	for key, val := range fieldValueMap {
		mapVal, ok := val.(map[string]interface{})
		if ok {
			var alreadyGetSubVal bool
			for subKey, subVal := range mapVal {
				if alreadyGetSubVal {
					break
				}
				switch subKey {
				case "$numberLong":  // NumberLong()
					fallthrough
				case "$binary":  // $type "05" ... // TODO 处理成原值
					fallthrough
				case "$code":  // Object()
					fallthrough
				case "t":  // i 1  // 时间戳（Timestamp() 秒）
					fallthrough
				case "$oid":  // ObjectId
					fallthrough
				case "$date":  // Date(), ISODate()
					alreadyGetSubVal = true
					val = subVal
				}
			}
		}
		result[key] = val
	}

	return result, nil
}
