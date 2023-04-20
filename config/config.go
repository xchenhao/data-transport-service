package config

import (
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"

	"github.com/xchenhao/data-transport-service/kafka"
	"github.com/xchenhao/data-transport-service/sql"
)

type Config struct {
	Kafka kafka.Config `yaml:"kafka"`
	DB sql.Config `yaml:"db"`
	MongoDBMapToSQLFile []string `yaml:"mongodb_map_to_sql"`
}

func LoadConfig(confPath string) (Config, error) {
	conf := Config{}

	configFile, err := ioutil.ReadFile(confPath)
	if err != nil {
		return conf, err
	}

	err = yaml.Unmarshal(configFile, &conf)
	if err != nil {
		return conf, err
	}

	return conf, nil
}