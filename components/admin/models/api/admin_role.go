// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/db"
)

func GetAdminRoles(adminID int) (roles []*models.Role) {
	qe := db.DB.Model(&models.AdminRole{}).Where("admin_id=?", adminID).Order("role_id asc").Select("role_id").QueryExpr()
	db.DB.Where("id in (?)", qe).Find(&roles)
	return
}

// 该操作本身应当设置权限控制
func AddAdminRole(adminID, roleID int) error {
	return db.DB.Create(&models.AdminRole{
		ID:      0,
		AdminID: adminID,
		RoleID:  roleID,
	}).Error
}

func DelAdminRole(adminID, roleID int) error {
	return db.DB.Where("admin_id=? AND role_id=?", adminID, roleID).Delete(&models.AdminRole{}).Error
}