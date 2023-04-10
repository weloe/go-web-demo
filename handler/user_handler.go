package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go-web-demo/component"
	"go-web-demo/handler/request"
	"go-web-demo/service"
)

func Login(c *gin.Context) {
	loginRequest := &request.Login{}
	err := c.ShouldBindBodyWith(loginRequest, binding.JSON)
	if err != nil {
		c.JSON(400, component.RestResponse{Code: -1, Message: " bind error"})
		return
	}
	token := service.Login(loginRequest)

	c.JSON(200, component.RestResponse{Code: 1, Data: token, Message: loginRequest.Username + " logged in successfully"})

}

func Logout(c *gin.Context) {
	token := c.Request.Header.Get("token")

	if token == "" {
		panic(fmt.Errorf("token error: token is nil"))
	}

	bytes, err := component.GlobalCache.Get(token)

	if err != nil {
		panic(fmt.Errorf("token error: failed to get username : %v", err))
	}

	username := string(bytes)
	// Authentication

	// Delete store current subject in cache
	err = component.GlobalCache.Delete(token)
	if err != nil {
		panic(fmt.Errorf("failed to delete current subject in cache: %w", err))
	}

	c.JSON(200, component.RestResponse{Code: 1, Data: token, Message: username + " logout in successfully"})
}

func Register(c *gin.Context) {
	register := &request.Register{}
	err := c.ShouldBindBodyWith(register, binding.JSON)
	if err != nil {
		c.JSON(400, component.RestResponse{Code: -1, Message: " bind error"})
		return
	}

	service.Register(register)

	c.JSON(200, component.RestResponse{Code: 1, Data: nil, Message: "register successfully"})
}
