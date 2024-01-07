package models

import (
	"errors"
	"github.com/jinzhu/gorm"
)

type Auth struct {
	ID       int    `json:"id" gorm:"type:int(11) AUTO_INCREMENT comment 'id';primary_key" json:"id"`
	Username string `json:"username" gorm:"type:varchar(100) comment '用户名';default:'';"`
	Password string `json:"password" gorm:"type:varchar(100) comment '密码';default:'';"`
}

// CheckAuth 检查用户信息
func CheckAuth(username, password string) (bool, error) {
	var auth Auth
	err := db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if auth.ID > 0 {
		return true, nil
	}

	return false, nil
}

// EditUser 编辑用户信息
func EditUser(username string, data interface{}) error {
	if err := db.Model(&Auth{}).Where("username = ? ", username).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

func AddUser(account string, password string) error {
	auth := Auth{
		Username: account,
		Password: password,
	}
	if err := db.Create(&auth).Error; err != nil {
		return err
	}
	return nil
}
