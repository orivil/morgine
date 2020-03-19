// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/cache"
)

func GetRoutes(routes []string) []*models.Route {
	return cache.RouteMux.GetRoutes(routes)
}

func AllRoutes() []*models.Route {
	return cache.RouteMux.AllRoutes()
}

func CheckRoutesExist(routes []string) (notExists []string) {
	return cache.RouteMux.CheckRoutesExist(routes)
}