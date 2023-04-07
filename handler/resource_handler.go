package handler

import (
	"github.com/gin-gonic/gin"
	"go-web-demo/component"
)

func ReadResource(c *gin.Context) {

	c.JSON(200, component.RestResponse{Code: 1, Message: "read resource successfully", Data: "resource"})
}

func WriteResource(c *gin.Context) {

	c.JSON(200, component.RestResponse{Code: 1, Message: "write resource successfully", Data: "resource"})
}
