package component

import (
	"fmt"
	"github.com/allegro/bigcache"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	_ "github.com/go-sql-driver/mysql"
	tokenGo "github.com/weloe/token-go"
	_ "github.com/weloe/token-go-extensions/redis-adapter"
	"go-web-demo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"time"
)

var (
	DB            *gorm.DB
	GlobalCache   *bigcache.BigCache
	Enforcer      *casbin.Enforcer
	TokenEnforcer *tokenGo.Enforcer
)

// CreateByConfig create components
func CreateByConfig() {

	ConnectDB()

	CreateLocalCache()

	CreateCasbinEnforcer()

	CreateTokenEnforcer()
}

func ConnectDB() {
	// connect to DB
	var err error
	dbConfig := config.Reader.ReadConfig().Mysql
	if dbConfig == nil {
		log.Fatalf(fmt.Sprintf("db config is nil"))
	}
	// config
	username := dbConfig.Username
	password := dbConfig.Password
	host := dbConfig.Host
	port := dbConfig.Port
	Dbname := dbConfig.Dbname
	timeout := dbConfig.TimeOut

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	log.Println("connect db url: " + dbUrl)
	DB, err = gorm.Open(mysql.Open(dbUrl), &gorm.Config{})

	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to connect to DB: %v", err))
	}
}

func CreateLocalCache() {
	var err error
	cacheConfig := config.Reader.ReadConfig().LocalCache
	if cacheConfig == nil {
		log.Fatalf(fmt.Sprintf("cache config is nil"))
	}
	// Initialize cache to store current user in cache.
	GlobalCache, err = bigcache.NewBigCache(bigcache.DefaultConfig(cacheConfig.ExpireTime * time.Second)) // Set expire time to 30 s
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to initialize cahce: %v", err))
	}
}

func CreateCasbinEnforcer() {
	var err error

	// casbin model
	config := config.Reader.ReadConfig().Casbin
	if config == nil {
		log.Fatalf(fmt.Sprintf("casbin config is nil"))
	}
	model := config.Model
	//Initialize casbin adapter
	adapter, _ := gormadapter.NewAdapterByDB(DB)

	// Load model configuration file and policy store adapter
	Enforcer, err = casbin.NewEnforcer(model, adapter)
	if err != nil {
		log.Fatalf(fmt.Sprintf("failed to create casbin enforcer: %v", err))
	}
	Enforcer.EnableEnforce(config.EnableEnforcer)
	Enforcer.EnableLog(config.EnableLog)

	//// Load policies from DB dynamically
	//err = Enforcer.LoadPolicy()
	//if err != nil {
	//	log.Fatalf(fmt.Sprintf("failed to load policy from DB: %v", err))
	//}
}

func CreateTokenEnforcer() {
	var err error
	// adapter, err := redisAdapter.NewAdapter("ip:host", "username", "pwd", 0)
	adapter := tokenGo.NewDefaultAdapter()
	if err != nil {
		log.Fatalf("NewRedisAdapter() failed: %v", err)
	}
	TokenEnforcer, err = tokenGo.NewEnforcer(adapter)
	if err != nil {
		log.Fatalf("NewDefaultEnforcer() failed: %v", err)
	}
}
