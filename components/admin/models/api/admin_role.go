// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"errors"
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/db"
)

func GetAdminRoles(parentID, adminID int) (roles []*models.Role, err error) {
	if parentID > 0 {
		if !IsSubAdmin(parentID, adminID) {
			return nil, errors.New("当前登录账号权限不足")
		}
	}
	qe := db.DB.Model(&models.AdminRole{}).Where("admin_id=?", adminID).Order("role_id asc").Select("role_id").QueryExpr()
	db.DB.Where("id in (?)", qe).Find(&roles)
	return
}

// 该操作本身应当设置权限控制
func AddAdminRole(parentID, adminID, roleID int) error {
	if parentID > 0 {
		if !IsSubAdmin(parentID, adminID) {
			return errors.New("当前登录账号权限不足")
		}
		if !isAdminHasRole(parentID, roleID) {
			return errors.New("当前登录账号没有该权限")
		}
	}
	return db.DB.Create(&models.AdminRole{
		ID:      0,
		AdminID: adminID,
		RoleID:  roleID,
	}).Error
}

func DelAdminRole(parentID, adminID, roleID int) error {
	if parentID > 0 {
		if !IsSubAdmin(parentID, adminID) {
			return errors.New("当前登录账号权限不足")
		}
	}
	return db.DB.Where("admin_id=? AND role_id=?", adminID, roleID).Delete(&models.AdminRole{}).Error
}

func isAdminHasRole(adminID, roleID int) bool {
	return IsIDExist(db.DB.Model(&models.AdminRole{}).Where("admin_id=? AND role_id=?", adminID, roleID))
}