// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/db"
	"github.com/orivil/morgine/components/admin/utils"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"strings"
)

var (
	ErrUsernameAlreadyExist = errors.New("用户名已注册")
	ErrNoDelAdminContainsSub = errors.New("不可删除包含子账号的账号")
	ErrLoginUsernameIncorrect = errors.New("用户名错误")
	ErrOnlyCanDelOwnSubAccount = errors.New("只能删除自己的后代账号")
	ErrMismatchedPassword = bcrypt.ErrMismatchedHashAndPassword
	ErrNeedSuperAdmin = errors.New("需要超级管理员才能操作")
	ErrOnlyCanCreateOwnSubAccount = errors.New("只能创建属于自己的子账号")
)

func RegisterSuperAdmin(username, nickname, password string) (admin *models.Admin, err error) {
	if IsIDExist(db.DB.Model(&models.Admin{})) {

	}
	admin = &models.Admin {
		Username:  username,
		Nickname:  nickname,
		Password:  password,
	}
	db.DB.Create(admin)
}

// 创建子管理员
func CreateSubAdmin(loginID, parentID int, admin *models.Admin) error {
	admin.ID = 0
	if IsIDExist(db.DB.Model(&models.Admin{}).Where("username=?", admin.Username)) {
		return ErrUsernameAlreadyExist
	} else {
		parent, err := getAccount(parentID)
		if err != nil {
			return err
		}
		if loginID != parentID && !strings.Contains(parent.Forefather, "|" + strconv.Itoa(loginID) + "|") {
			return ErrOnlyCanCreateOwnSubAccount
		}
		if !parent.Super.IsTrue() {
			return ErrNeedSuperAdmin
		}
		var pass []byte
		pass, err = bcrypt.GenerateFromPassword([]byte(admin.Password), utils.RandomInt(bcrypt.MinCost, bcrypt.MaxCost))
		if err != nil {
			return err
		}
		admin.Password = string(pass)
		admin.ParentID = parentID
		admin.Forefather = parent.Forefather + "|" + strconv.Itoa(parentID)
		return db.DB.Create(admin).Error
	}
}

func GetAdminByID(id int) *models.Admin {
	admin := &models.Admin{}
	db.DB.Where("id=?", id).First(admin)
	if admin.ID > 0 {
		admin.Password = ""
		return admin
	} else {
		return nil
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
		a1.Password = ""
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

func DelSubAccount(parentID, childID int) error {
	arg := "%|" + strconv.Itoa(childID) + "|%"
	if IsIDExist(db.DB.Where("forefather LIKE ?", arg)) { // 检测被删除的账号是否存在子孙账号
		return ErrNoDelAdminContainsSub
	} else {
		arg = "%|" + strconv.Itoa(parentID) + "|%"
		anum := db.DB.Where("id=? AND forefather LIKE ?", childID, parentID).Delete(&models.Admin{}).RowsAffected
		if anum == 1 {
			return nil
		} else if anum == 0 {
			return ErrOnlyCanDelOwnSubAccount
		} else {
			return errors.New("invalid operation")
		}
	}
}

func LoginAdmin(username, password string) (admin *models.Admin, err error) {
	admin = &models.Admin{}
	err = db.DB.Where("username=?", username).First(admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrLoginUsernameIncorrect
		} else {
			return nil, err
		}
	} else {
		err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
		if err != nil {
			if err == bcrypt.ErrMismatchedHashAndPassword {
				return nil, ErrMismatchedPassword
			} else {
				return nil, err
			}
		}
		admin.Password = ""
		return admin, err
	}
}

func UpdateAdminPassword(parentID, subID int, password string) error {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), utils.RandomInt(bcrypt.MinCost, bcrypt.MaxCost))
	if err != nil {
		return err
	}
	var affected int64
	if subID > 0 {
		affected = db.DB.Model(&models.Admin{}).Where("id=? AND forefather LIKE ?", subID, "%|" + strconv.Itoa(parentID) + "|%").UpdateColumn("password", string(pass)).RowsAffected
	} else {
		affected = db.DB.Model(&models.Admin{}).Where("id=?", parentID).UpdateColumn("password", string(pass)).RowsAffected
	}
	if affected == 1 {
		return nil
	} else {
		return errors.New("failed")
	}
}

func UpdateAdminInfo(parentID, subID int, info *models.Admin) error {
	var affected int64
	if subID > 0 {
		arg := "%|" + strconv.Itoa(parentID) + "|%"
		affected = db.DB.Model(&models.Admin{}).Where("id=? AND forefather LIKE ?", subID, arg).Updates(info).RowsAffected
	} else {
		affected = db.DB.Model(&models.Admin{}).Where("id=?", parentID).Updates(info).RowsAffected
	}
	if affected == 1 {
		return nil
	} else {
		return errors.New("failed")
	}
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

func IsSubAdmin(parentID, childID int) bool {
	return IsIDExist(db.DB.Model(&models.Admin{}).Where("id=? AND forefather LIKE ?", childID, "%|" + strconv.Itoa(parentID) + "|%"))
}