package sql

import (
	"fmt"
	"time"

	// _ "database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Config struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
	Database string `yaml:"database"`
	User string `yaml:"user"`
	Password string `yaml:"password"`
	MaxLifeTime string `yaml:"max_life_time"`
	MaxIdleConns int    `yaml:"max_idle_conns"`
	MaxOpenConns int    `yaml:"max_open_conns"`
	Charset string `yaml:"charset"`
	Debug bool `yaml:"debug"`
}

var defaultMaxLifeTime = "10s"
var defaultCharset = "utf8mb4,utf8"

func InitDB(conf Config) (*gorm.DB, error) {
	confMaxLifeTime := conf.MaxLifeTime
	if conf.MaxLifeTime == "" {
		confMaxLifeTime = "10s"
	}
	maxLifeTime, err := time.ParseDuration(confMaxLifeTime)
	if err != nil {
		return nil, err
	}
	charset := conf.Charset
	if charset == "" {
		charset = defaultCharset
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		conf.User, conf.Password, conf.Host,
		conf.Port, conf.Database, charset, true, "Local")
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// db.BlockGlobalUpdate(true)
	if conf.Debug {
		db = db.Debug()
	}
	sqldb := db.DB()
	err = sqldb.Ping()
	if err != nil {
		db.Close()
		return nil, err
	}
	sqldb.SetMaxIdleConns(conf.MaxIdleConns)
	// Fix db invalid connection after EOF
	sqldb.SetConnMaxLifetime(maxLifeTime)
	if conf.MaxOpenConns != 0 {
		sqldb.SetMaxOpenConns(conf.MaxOpenConns)
	}

	return db, nil
}
