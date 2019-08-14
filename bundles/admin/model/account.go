// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin_model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/orivil/morgine/bundles/utils/sql"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameRegistered = errors.New("该用户名已注册")
	ErrUserNotRegistered  = errors.New("用户不存在")
	ErrPasswordIncorrect  = errors.New("密码错误")
	ErrUsernameIncorrect  = errors.New("用户名不存在")
)

type Account struct {
	ID       int
	Super sql.Boolean
	Username string `gorm:"index"`
	Password string
}

func CreateAdmin(username, password string, super bool) error {
	var a = &Account{}
	DB.Model(a).Where("username=?", a.Username).Select("id").First(a)
	if a.ID > 0 {
		return ErrUsernameRegistered
	}
	password, err := hashPassword(password)
	if err != nil {
		return err
	}
	return DB.Create(&Account{
		Username: username,
		Password: password,
		Super: sql.GetSqlBoolean(super),
	}).Error
}

func UpdatePassword(username, oldPassword, newPassword string) error {
	var exist = &Account{}
	DB.Model(exist).Where("username=?", username).First(exist)
	if exist.ID == 0 {
		return ErrUserNotRegistered
	}
	err := bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(oldPassword))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return ErrPasswordIncorrect
		} else {
			return err
		}
	}
	newPassword, err = hashPassword(newPassword)
	if err != nil {
		return err
	}
	return DB.Model(exist).Where("username=?", username).UpdateColumn("password", newPassword).Error
}

func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	} else {
		return string(pw), nil
	}
}

func SignIn(username, password string) (id int, err error) {
	exist := &Account{}
	err = DB.Where("username=?", username).Select("id").First(exist).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, ErrUsernameIncorrect
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return 0, ErrPasswordIncorrect
		} else {
			return 0, err
		}
	}
	return exist.ID, nil
}

func GetAdmin(id int) (*Account, error) {
	admin := &Account{}
	err := DB.Where("id=?", id).First(admin).Error
	if err != nil {
		return nil, err
	} else {
		return admin, nil
	}
}
