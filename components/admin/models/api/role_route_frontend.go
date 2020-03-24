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

func GetFrontendRolesRoutes(roleID []int) (rs []*models.RoleRouteFrontend) {
	db.DB.Where("role_id IN ?", roleID).Order("id asc").Find(&rs)
	return
}

func CreateFrontendRoleRoute(roleID int, route string) (*models.RoleRouteFrontend, error) {
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
}

func DelFrontendRoleRoute(roleID int, route string) error {
	num := db.DB.Where("role_id=? AND route=?", roleID, route).Delete(&models.RoleRouteFrontend{}).RowsAffected
	if num > 0 {
		return nil
	} else {
		return errors.New("failed")
	}
}

func GetAdminFrontendRoleRoutes(adminID int) (routes []*models.RoleRouteFrontend) {
	qexp := db.DB.Model(&models.AdminRole{}).Where("admin_id=?", adminID).Select("role_id")
	db.DB.Where("role_id in ?", qexp).Order("id asc").Find(&routes)
	return routes
}