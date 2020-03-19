// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"fmt"
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/db"
)

func GetRoleRoutes(roleID int) (routes []*models.Route, err error) {
	var rs []string
	rs, err = selectRouteRoutes(roleID)
	if err != nil {
		return nil, err 
	} else {
		routes = GetRoutes(rs)
		return routes, nil
	}
}

func GetRoleDeletedRoutes(roleID int) (routes []string, err error) {
	routes, err = selectRouteRoutes(roleID)
	if err != nil {
		return nil, err
	} else {
		routes = CheckRoutesExist(routes)
		return routes, nil
	}
}

func AddRoleRoute(roleID int, route string) error {
	notExist := CheckRoutesExist([]string{route})
	if len(notExist) > 0 {
		return fmt.Errorf("route: '%s' is not exist", route)
	} else {
		if IsIDExist(db.DB.Model(&models.RoleRoute{}).Where("role_id=? AND route=?", roleID, route)) {
			return fmt.Errorf("route: '%s' already exist")
		} else {
			return db.DB.Create(&models.RoleRoute{
				ID:     0,
				RoleID: roleID,
				Route:  route,
			}).Error
		}
	}
}

func DelRoleRoute(roleID int, route string) error {
	return  db.DB.Where("role_id=? AND route=?", route).Delete(&models.RoleRoute{}).Error
}

func selectRouteRoutes(roleID int) (routes []string, err error) {
	err = db.DB.Model(&models.RoleRoute{}).Where("role_id=?", roleID).Order("id desc").Pluck("route", &routes).Error
	return routes, err
}