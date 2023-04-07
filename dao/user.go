package dao

import "go-web-demo/component"

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Username string
	Password string
	Email    string
	Phone    string
}

func GetByUsername(username string) *User {
	res := new(User)
	component.DB.Model(&User{}).Where("username = ?", username).First(res)
	return res
}
