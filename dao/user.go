package dao

import "go-web-demo/component"

type User struct {
	Id       int64 `gorm:"primaryKey"`
	Username string
	Password string
	Email    string
	Phone    string
}

func (u *User) TableName() string {
	return "user"
}

func GetByUsername(username string) *User {
	res := new(User)
	component.DB.Model(&User{}).Where("username = ?", username).First(res)
	return res
}

func Insert(username string, password string) (int64, error, int64) {
	user := &User{Username: username, Password: password}
	res := component.DB.Create(&user)

	return user.Id, res.Error, res.RowsAffected
}
