package main

import (
	"fmt"
	gormadapter "github.com/casbin/gorm-adapter/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go-web-demo/component"
	"go-web-demo/config"
	"go-web-demo/handler"
	"go-web-demo/middleware"
	"log"
)

var (
	router *gin.Engine
)

func init() {
	// Initialize confif yaml
	reader := config.Reader.ReadConfig()
	// casbin model
	model := reader.Casbin.Model

	//Initialize components: mysql locaCache
	component.CreateByConfig()

	//Initialize casbin adapter
	adapter, _ := gormadapter.NewAdapterByDB(component.DB)

	// Initialize gin router
	router = gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))
	router.Use(middleware.Recover)

	user := router.Group("/user")
	{
		user.POST("/login", handler.Login)
		user.POST("/logout", handler.Logout)
	}

	resource := router.Group("/api")
	{
		resource.GET("/resource", middleware.Authorize("resource", "read", adapter, model), handler.ReadResource)
		resource.POST("/resource", middleware.Authorize("resource", "write", adapter, model), handler.WriteResource)
	}

}

func main() {
	// Start
	defer component.DB.Close()
	port := "8081"
	err := router.Run(":" + port)
	if err != nil {
		panic(fmt.Sprintf("failed to start gin engine: %v", err))
	}
	log.Println("application is now running...")
}
