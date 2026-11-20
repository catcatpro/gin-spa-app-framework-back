package models

import (
	"errors"
	"gin-spa-app-framework-back/common"
	"gin-spa-app-framework-back/utils"
)

type User struct {
	Basic
	Username string `json:"username" gorm:"NOT NULL"`
	Password string `json:"password" gorm:"NOT NULL"`
}

// func (User) TableName() string {
// 	return "user"
// }

// 用户登录
func (u *User) UserLogin() error {
	var password string = u.Password
	res := common.COM_DB.Where("username = ?", u.Username).First(&u)
	// fmt.Println("debug: ", u.Password, findUser.Password)
	verify_res := utils.CheckPasswordHash(password, u.Password)
	if !verify_res {
		return errors.New("Username or Password incorrect!")
	}

	return res.Error
}

// 新建用户
func (u *User) CreateUser() (uint, error) {
	var saveData User
	saveData = *u
	saveData.Password, _ = utils.HashPassword(saveData.Password)
	res := common.COM_DB.Create(&saveData)
	return saveData.ID, res.Error
}

// 通过token获取用户信息
func (u *User) GetUserInfoByToken(token string) error {
	j := utils.NewJWT()
	user, err := j.ParseToken(token, false)
	if err != nil {
		return err
	}
	res := common.COM_DB.Where("id = ?", user.Id).Select("id, username").First(&u)
	return res.Error
}

// 更改系统用户密码
func (u *User) UpdateUserInfo() error {
	password, _ := utils.HashPassword(u.Password)
	res := common.COM_DB.Model(&u).Where("id = ?", u.ID).Update("password", password)
	return res.Error
}
