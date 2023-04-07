package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go-web-demo/component"
	"go-web-demo/handler/request"
)

func Login(c *gin.Context) {
	loginRequest := &request.Login{}
	err := c.ShouldBindBodyWith(loginRequest, binding.JSON)
	if err != nil {
		c.JSON(400, component.RestResponse{Code: -1, Message: " bind error"})
		return
	}
	password := loginRequest.Password
	username := loginRequest.Username

	// Authentication 123 for test
	if password != "123" {
		panic(fmt.Errorf(username + " logged error : password error"))
	}

	// Generate random session id
	u, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("failed to generate UUID: %w", err))
	}
	// sprintf
	token := fmt.Sprintf("%s-%s", u.String(), "token")
	// Store current subject in cache
	err = component.GlobalCache.Set(token, []byte(username))
	if err != nil {
		panic(fmt.Errorf("failed to store current subject in cache: %w", err))
	}
	// Send cache key back to client cookie
	//c.SetCookie("current_subject", token, 30*60, "/resource", "", false, true)

	c.JSON(200, component.RestResponse{Code: 1, Data: token, Message: username + " logged in successfully"})

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
