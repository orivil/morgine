// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"errors"
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/db"
	"github.com/orivil/morgine/utils/sql"
	"strconv"
	"time"
)

var (
	ErrUsernameAlreadyExist = errors.New("用户名已注册")
)

// 创建子管理员
func CreateSubAdmin(parentID int, admin *models.Admin) error {
	admin.ID = 0
	if IsIDExist(db.DB.Model(&models.Admin{}).Where("username=?", admin.Username)) {
		return ErrUsernameAlreadyExist
	} else {
		if admin.Super.IsTrue() {
			var super []sql.Boolean
			db.DB.Model(&models.Admin{}).Where("id=?").Limit(1).Pluck("super", &super)
			if len(super) != 1 {
				return errors.New("parent admin not exist")
			} else {
				if !super[0].IsTrue() {
					return errors.New("非超级管理员不可创建超级管理员")
				}
			}
		}
		parent, err := getAccount(parentID)
		if err != nil {
			return err
		}
		admin.ParentID = parentID
		admin.Level = parent.Level + 1
		admin.Forefather = parent.Forefather + "," +
		return db.DB.Create(admin).Error
	}
}

type Account struct {
	*models.Admin
	Subs []*Account
}

// 获得所有子管理员列表
func GetSubAdmins(parentID int) (accounts []*Account) {
	// 找到所有子账号
	var admins []*models.Admin
	arg := "%|" + strconv.Itoa(parentID) + "|%"
	db.DB.Where("forefather LIKE ?", arg).Order("id desc").Find(&admins)
	for _, a1 := range admins {
		account := &Account {
			Admin: a1,
		}
		for _, a2 := range admins {
			if a2.ParentID == a1.ID {
				account.Subs = append(account.Subs, &Account{
					Admin: a2,
				})
			}
		}
		if a1.ParentID == parentID {
			accounts = append(accounts, account)
		}
	}
	return accounts
}

func DelSubAccount(parentID, childID int) {

}

func LoginAdmin(account, password string) error {
	db.DB.Where("")
}

func getAccount(id int) (admin *models.Admin, err error) {
	admin = &models.Admin{}
	err = db.DB.Where("id=?", id).First(admin).Error
	if err != nil {
		return nil, err
	} else {
		return admin, nil
	}
}