// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"errors"
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/db"
)

var (
	ErrNoDelRoleUsed = errors.New("不可删除已被使用的权限")
	ErrNoDelParentRole = errors.New("不可删除父辈角色")
)

func GetAllRoles() (roles []*models.Role) {
	db.DB.Order("id asc").Find(&roles)
	return
}

func CreateRole(role *models.Role) error {
	return db.DB.Create(role).Error
}

func UpdateRole(id int, role *models.Role) error {
	role.ID = 0
	return db.DB.Model(role).Where("id=?", id).Updates(role).Error
}

// 删除角色
//
// 如果强制删除角色，则关联的管理员角色配置，角色路由配置也将同时被删除，如果不强制删除，
// 则会检查管理员角色配置及角色路由配置，如果找到该角色已有配置则将会停止删除角色并返回
// ErrNoDelRoleUsed 错误
//
//任何情况下都不可删除父辈角色（即包含子角色），否则会报 ErrNoDelParentRole 错误
func DeleteRole(id int, force bool) error {
	if IsIDExist(db.DB.Model(&models.Role{}).Where("parent_id=?", id)) {
		return ErrNoDelParentRole
	}
	if force {
		err := db.DB.Where("role_id=?", id).Delete(&models.RoleRoute{}).Error
		if err != nil {
			return err
		}
		err = db.DB.Where("role_id=?", id).Delete(&models.AdminRole{}).Error
		if err != nil {
			return err
		}
		return db.DB.Where("id=?", id).Delete(&models.Role{}).Error
	} else {
		if !IsIDExist(db.DB.Model(&models.RoleRoute{}).Where("role_id=?")) &&
			!IsIDExist(db.DB.Model(&models.AdminRole{}).Where("role_id=?")){
			return db.DB.Where("id=?", id).Delete(&models.Role{}).Error
		} else {
			return ErrNoDelRoleUsed
		}
	}
}