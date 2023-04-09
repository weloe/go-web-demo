package component

import (
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"go-web-demo/config"
	"log"
	"time"
)

var (
	DB          *gorm.DB
	GlobalCache *bigcache.BigCache
	Enforcer    *casbin.Enforcer
)

// CreateByConfig create components
func CreateByConfig() {

	ConnectDB()

	CreateLocalCache()

	CreateCasbinEnforcer()
}

func ConnectDB() {
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
	log.Println("connect db url: " + dbUrl)
	DB, err = gorm.Open("mysql", dbUrl)
	DB.SingularTable(true)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to connect to DB: %v", err))
	}
}

func CreateLocalCache() {
	var err error
	cacheConfig := config.Reader.LocalCache
	// Initialize cache to store current user in cache.
	GlobalCache, err = bigcache.NewBigCache(bigcache.DefaultConfig(cacheConfig.ExpireTime * time.Second)) // Set expire time to 30 s
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to initialize cahce: %v", err))
	}
}

func CreateCasbinEnforcer() {
	var err error

	// casbin model
	model := config.Reader.Casbin.Model
	//Initialize casbin adapter
	adapter, _ := gormadapter.NewAdapterByDB(DB)

	// Load model configuration file and policy store adapter
	Enforcer, err = casbin.NewEnforcer(model, adapter)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}
	
	//// Load policies from DB dynamically
	//err = Enforcer.LoadPolicy()
	//if err != nil {
	//	log.Fatalf(fmt.Sprintf("failed to load policy from DB: %v", err))
	//}
}
