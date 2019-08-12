// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package model

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/orivil/morgine/bundles/admin/db"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUsernameRegistered = errors.New("该用户名已注册")
	ErrUserNotRegistered  = errors.New("用户不存在")
	ErrPasswordIncorrect  = errors.New("密码错误")
	ErrUsernameIncorrect  = errors.New("用户名错误")
)

type Admin struct {
	ID       int
	Username string `gorm:"index"`
	Password string
}

func CreateAdmin(username, password string) error {
	var a = &Admin{}
	db.GORM.Model(a).Where("username=?", a.Username).Select("id").First(a)
	if a.ID > 0 {
		return ErrUsernameRegistered
	}
	password, err := hashPassword(password)
	if err != nil {
		return err
	}
	return db.GORM.Create(&Admin{
		Username: username,
		Password: password,
	}).Error
}

func UpdatePassword(username, oldPassword, newPassword string) error {
	var exist = &Admin{}
	db.GORM.Model(exist).Where("username=?", username).First(exist)
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
	return db.GORM.Model(exist).Where("username=?", username).UpdateColumn("password", newPassword).Error
}

func hashPassword(password string) (string, error) {
	pw, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	} else {
		return string(pw), nil
	}
}

func SignIn(username, password string) (a *Admin, err error) {
	exist := &Admin{}
	err = db.GORM.Where("username=?", username).First(exist).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUsernameIncorrect
		} else {
			return nil, err
		}
	}
	err = bcrypt.CompareHashAndPassword([]byte(exist.Password), []byte(password))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return nil, ErrPasswordIncorrect
		} else {
			return nil, err
		}
	}
	return exist, nil
}

func GetAdmin(id int) (*Admin, error) {
	admin := &Admin{}
	err := db.GORM.Where("id=?", id).First(admin).Error
	if err != nil {
		return nil, err
	} else {
		return admin, nil
	}
}
