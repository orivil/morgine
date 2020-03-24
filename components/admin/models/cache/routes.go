// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cache

type Route struct {
	Title string
	Method string
	Path string
}

var FrontendRouteMux = &routeMux {
	ls: []*Route{},
	ms: make(map[string]*Route, 50),
}

type routeMux struct {
	ls []*Route
	ms map[string]*Route
}

func (rm *routeMux) Init(title, method, path string) {
	route := &Route {
		Title:  title,
		Path:   path,
	}
	rm.ls = append(rm.ls, route)
	rm.ms[method + path] = route
}

func (rm *routeMux) AllRoutes() []*Route {
	return rm.ls
}

func (rm *routeMux) GetRoutes(routes []string) (rs []*Route) {
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
