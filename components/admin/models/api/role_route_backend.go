// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"errors"
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/db"
)
var ErrBackendRoleRouteAlreadyExist = errors.New("backend role route already exist")

func GetBackendRoleRoutes(roleID int) (routes []*models.RoleRouteBackend) {
	db.DB.Where("role_id = ?", roleID).Order("id asc").Find(&routes)
	return
}

func CreateBackendRoleRoutes(adminID int, roleID int, route string) (*models.RoleRouteBackend, error) {
	if IsSuperAdmin(adminID) {
		if IsIDExist(db.DB.Model(&models.RoleRouteBackend{}).Where("role_id=? AND route=?", roleID, route)) {
			return nil, ErrBackendRoleRouteAlreadyExist
		}
		m := &models.RoleRouteBackend{
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
	} else {
		return nil, ErrNeedSuperAdmin
	}
}

func DelBackendRoleRoute(adminID, roleID int, route string) error {
	if IsSuperAdmin(adminID) {
		num := db.DB.Where("role_id=? AND route=?", roleID, route).Delete(&models.RoleRouteBackend{}).RowsAffected
		if num > 0 {
			return nil
		} else {
			return errors.New("failed")
		}
	} else {
		return ErrNeedSuperAdmin
	}
}

func GetAdminBackendRoleRoutes(adminID int) (routes []*models.RoleRouteBackend) {
	qexp := db.DB.Model(&models.AdminRole{}).Where("admin_id=?", adminID).Select("role_id")
	db.DB.Where("role_id in ?", qexp).Order("id asc").Find(&routes)
	return routes
}