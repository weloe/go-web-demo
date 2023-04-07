package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync"
	"time"
)

type Config struct {
	Server     *Server
	Mysql      *DB
	LocalCache *LocalCache
	Casbin     *Casbin
}

type Server struct {
	Port int64
}

type DB struct {
	Username string
	Password string
	Host     string
	Port     int64
	Dbname   string
	TimeOut  string
}

type LocalCache struct {
	ExpireTime time.Duration
}

type Casbin struct {
	Model string
}

var (
	once   sync.Once
	Reader = new(Config)
)

func (config *Config) ReadConfig() *Config {
	once.Do(func() {
		viper.SetConfigName("config")   // filename
		viper.SetConfigType("yaml")     // filename extension : yaml | json |
		viper.AddConfigPath("./config") // workspace dir : ./
		var err error
		err = viper.ReadInConfig() // read config
		if err != nil {            // handler err
			log.Fatalf(fmt.Sprintf("Fatal error config file: %s \n", err))
		}
		err = viper.Unmarshal(config)
		if err != nil {
			log.Fatalf(fmt.Sprintf("Fatal error viper unmarshal config: %s \n", err))
		}
	})
	return Reader
}
