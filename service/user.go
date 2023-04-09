package service

import (
	"fmt"
	"github.com/google/uuid"
	"go-web-demo/component"
	"go-web-demo/dao"
	"go-web-demo/handler/request"
)

func Login(loginRequest *request.Login) string {
	password := loginRequest.Password
	username := loginRequest.Username

	// Authentication
	user := dao.GetByUsername(username)
	if password != user.Password {
		panic(fmt.Errorf(username + " logged error : password error"))
	}

	// Generate random uuid token
	u, err := uuid.NewRandom()
	if err != nil {
		panic(fmt.Errorf("failed to generate UUID: %w", err))
	}
	// Sprintf token
	token := fmt.Sprintf("%s-%s", u.String(), "token")
	// Store current subject in cache
	err = component.GlobalCache.Set(token, []byte(username))
	if err != nil {
		panic(fmt.Errorf("failed to store current subject in cache: %w", err))
	}
	// Send cache key back to client cookie
	//c.SetCookie("current_subject", token, 30*60, "/resource", "", false, true)
	return token
}
