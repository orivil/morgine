// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cache

import (
	"github.com/orivil/morgine/components/admin/models"
)

var RouteMux = &routeMux{
	ls: []*models.Route{},
	ms: make(map[string]*models.Route, 50),
}

type routeMux struct {
	ls []*models.Route
	ms map[string]*models.Route
}

func (rm *routeMux) Init(title, method, path string) {
	route := &models.Route {
		Title:  title,
		Method: method,
		Path:   path,
	}
	rm.ls = append(rm.ls, route)
	rm.ms[method + path] = route
}

func (rm *routeMux) AllRoutes() []*models.Route {
	return rm.ls
}

func (rm *routeMux) GetRoutes(routes []string) (rs []*models.Route) {
	for _, route := range routes {
		if r, ok := rm.ms[route]; ok {
			rs = append(rs, r)
		}
	}
	return rs
}

// 检测路由是否存在并返回不存在的路由
func (rm *routeMux) CheckRoutesExist(routes []string) (notExists []string) {
	for _, route := range routes {
		if _, ok := rm.ms[route]; !ok {
			notExists = append(notExists, route)
		}
	}
	return
}
