package kafka

import (
	"strings"

	"github.com/segmentio/kafka-go"
)

type Config struct {
	BootstrapServers string `yaml:"bootstrap_server"`  // comma-seperated
	Topics string `yaml:"topics"`  // comma-seperated
	GroupId string `yaml:"group_id"`
}

func Reader(conf Config) *kafka.Reader {
	brokers := strings.Split(conf.BootstrapServers, ",")
	topics := strings.Split(conf.Topics, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                brokers,
		GroupID:                conf.GroupId,
		// Topic:                  topic,
		GroupTopics:            topics,
		MinBytes:               10e3,  // 10KB
		MaxBytes:               10e6,  // 10MB
	})
}

