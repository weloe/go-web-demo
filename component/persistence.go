package component

import (
	"fmt"
	"github.com/allegro/bigcache"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-web-demo/config"
	"log"
	"time"
)

var (
	DB          *gorm.DB
	GlobalCache *bigcache.BigCache
)

func CreateByConfig() {

	connectDB()

	createLocalCache()

}

func connectDB() {
	// connect to DB
	var err error
	dbConfig := config.Reader.Mysql
	// config
	username := dbConfig.Username
	password := dbConfig.Password
	host := dbConfig.Host
	port := dbConfig.Port
	Dbname := dbConfig.Dbname
	timeout := dbConfig.TimeOut

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	log.Println(dbUrl)
	DB, err = gorm.Open("mysql", dbUrl)
	DB.SingularTable(true)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to connect to DB: %v", err))
	}
}

func createLocalCache() {
	var err error
	cacheConfig := config.Reader.LocalCache
	// Initialize cache to store current user in cache.
	GlobalCache, err = bigcache.NewBigCache(bigcache.DefaultConfig(cacheConfig.ExpireTime * time.Second)) // Set expire time to 30 s
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to initialize cahce: %v", err))
	}
}
