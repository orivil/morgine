// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin_model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameRegistered = errors.New("该用户名已注册")
	ErrUserNotRegistered  = errors.New("用户不存在")
	ErrPasswordIncorrect  = errors.New("密码错误")
	ErrUsernameIncorrect  = errors.New("用户名不存在")
)

type Admin struct {
	ID       int
	Username string `gorm:"unique_index"`
	Password string
	RoleID int `gorm:"index"`
}

func CountAdmins() (total int, err error) {
	err = DB.Model(&Admin{}).Count(&total).Error
	return
}

func CreateAdmin(username, password string) error {
	var a = &Admin{}
	DB.Model(a).Where("username=?", a.Username).Select("id").First(a)
	if a.ID > 0 {
		return ErrUsernameRegistered
	}
	password, err := HashPassword(password)
	if err != nil {
		return err
	}
	return DB.Create(&Admin{
		Username: username,
		Password: password,
	}).Error
}

func UpdatePassword(loginID int, username, newPassword string) error {
	var exist = &Admin{}
	DB.Model(exist).Where("id=? AND username=?", loginID, username).First(exist)
	if exist.ID == 0 {
		return ErrUserNotRegistered
	}
	newPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	return DB.Model(exist).Where("username=?", username).UpdateColumn("password", newPassword).Error
}

func HashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	} else {
		return string(pw), nil
	}
}

func SignIn(username, password string) (admin *Admin, err error) {
	admin = &Admin{}
	err = DB.Where("username=?", username).First(admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUsernameIncorrect
		} else {
			return nil, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return 0, ErrPasswordIncorrect
		} else {
			return 0, err
		}
	}
	return admin, nil
}

func GetAdmin(id int) (*Admin, error) {
	admin := &Admin{}
	err := DB.Where("id=?", id).First(admin).Error
	if err != nil {
		return nil, err
	} else {
		return admin, nil
	}
}
