package main

import (
	"fmt"
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
	//Initialize components from config yaml: mysql locaCache casbin
	component.CreateByConfig()

	// Initialize gin engine
	router = gin.Default()

	// Initialize gin middleware
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	router.Use(cors.New(corsConfig))
	router.Use(middleware.Recover)

	// Initialize gin router
	user := router.Group("/user")
	{
		user.POST("/login", handler.Login)
		user.POST("/logout", handler.Logout)
		user.POST("/register", handler.Register)
	}

	test := router.Group("/test")
	{
		test.POST("/t1", func(context *gin.Context) {

		})
		test.POST("/t2", func(context *gin.Context) {

		})
	}

	resource := router.Group("/api")
	{
		resource.Use(middleware.DefaultAuthorize("user::resource", "read-write"))
		resource.GET("/resource", handler.ReadResource)
		resource.POST("/resource", handler.WriteResource)
	}

}

func main() {
	// Start
	port := config.Reader.Server.Port
	err := router.Run(":" + port)
	if err != nil {
		panic(fmt.Sprintf("failed to start gin engine: %v", err))
	}
	log.Println("application is now running...")
}
