package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	tokenGo "github.com/weloe/token-go"
	"go-web-demo/component"
	"go-web-demo/handler/request"
	"go-web-demo/service"
	"net/http"
)

func Login(c *gin.Context) {
	loginRequest := &request.Login{}
	err := c.ShouldBindBodyWith(loginRequest, binding.JSON)
	if err != nil {
		panic(fmt.Errorf("request body bind error: %v", err))
	}

	token := service.Login(loginRequest, c)

	c.JSON(http.StatusOK, component.RestResponse{Code: 1, Data: token, Message: loginRequest.Username + " logged in successfully"})

}

func Logout(c *gin.Context) {
	// logout
	err := component.TokenEnforcer.Logout(tokenGo.NewHttpContext(c.Request, c.Writer))

	if err != nil {
		panic(fmt.Errorf("failed to Logout: %v", err))
	}

	c.JSON(http.StatusOK, component.RestResponse{Code: 1, Message: " logout in successfully"})
}

func Register(c *gin.Context) {
	register := &request.Register{}
	err := c.ShouldBindBodyWith(register, binding.JSON)
	if err != nil {
		c.JSON(400, component.RestResponse{Code: -1, Message: " bind error"})
		return
	}

	service.Register(register)

	c.JSON(http.StatusOK, component.RestResponse{Code: 1, Data: nil, Message: "register successfully"})
}
