// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"errors"
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/db"
)

var ErrFrontendRoleRouteAlreadyExist = errors.New("frontend role route already exist")

func GetFrontendRolesRoutes(adminID, roleID int) (rs []*models.RoleRouteFrontend, err error) {
	if IsSuperAdmin(adminID) {
		db.DB.Where("role_id = ?", roleID).Order("id asc").Find(&rs)
		return rs, nil
	} else {
		return nil, ErrNeedSuperAdmin
	}
}

func CreateFrontendRoleRoute(adminID, roleID int, route string) (*models.RoleRouteFrontend, error) {
	if IsSuperAdmin(adminID) {
		if IsIDExist(db.DB.Where("role_id=? AND route=?", roleID, route)) {
			return nil, ErrFrontendRoleRouteAlreadyExist
		} else {
			m := &models.RoleRouteFrontend{
				ID:     0,
				RoleID: roleID,
				Route:  route,
			}
			err := db.DB.Create(m).Error
			if err != nil {
				return nil, err
			} else {
				return m, nil
			}
		}
	} else {
		return nil, ErrNeedSuperAdmin
	}
}

func DelFrontendRoleRoute(adminID, roleID int, route string) error {
	if IsSuperAdmin(adminID) {
		num := db.DB.Where("role_id=? AND route=?", roleID, route).Delete(&models.RoleRouteFrontend{}).RowsAffected
		if num > 0 {
			return nil
		} else {
			return errors.New("failed")
		}
	} else {
		return ErrNeedSuperAdmin
	}
}

// TODO 测试是否会去重
func GetAdminFrontendRoleRoutes(adminID int) (routes []*models.RoleRouteFrontend) {
	qexp := db.DB.Model(&models.AdminRole{}).Where("admin_id=?", adminID).Select("role_id")
	db.DB.Where("role_id in ?", qexp).Order("id asc").Find(&routes)
	return routes
}