package service

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
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

func Register(register *request.Register) {
	var err error
	e := component.Enforcer
	err = e.GetAdapter().(*gormadapter.Adapter).Transaction(e, func(copyEnforcer casbin.IEnforcer) error {
		// Insert to table
		db := copyEnforcer.GetAdapter().(*gormadapter.Adapter).GetDb()
		res := db.Exec("insert into user (username,password) values(?,?)", register.Username, register.Password)

		//User has Username and Password
		//res := db.Table("user").Create(&User{
		//	Username: register.Username,
		//	Password: register.Password,
		//})

		if err != nil || res.RowsAffected < 1 {
			return fmt.Errorf("insert error: %w", err)
		}

		_, err = copyEnforcer.AddRoleForUser(register.Username, "role::user")
		if err != nil {
			return fmt.Errorf("add plocy error: %w", err)
		}

		return nil
	})

	if err != nil {
		panic(err)
	}

}

