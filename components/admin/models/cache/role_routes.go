// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cache

type RoleRoutesStorage interface {
	SetRoleRoute(role, route string)
	GetRoleRoutes(role string) (routes []string)
	DelRoleRoute(role, route string)
	MatchRoleRoute(role, route string) bool
}

type RoleRouteJsonStorage struct {
	filename string
	data map[string][]string
	mp map[string]map[string]struct{}
}

func (r *RoleRouteJsonStorage) SetRoleRoute(role, route string) {
	routes, _ := r.data[role]
	for _, rout := range routes {
		if rout == route {
			return
		}
	}
	r.data[role] = append(r.data[role], route)
	mp, ok := r.mp[role]
	if !ok {
		mp = make(map[string]struct{}, 1)
	}
	mp[route] = struct{}{}
	r.mp[role] = mp
}

func (r *RoleRouteJsonStorage) GetRoleRoutes(role string) (routes []string) {
	return r.data[role]
}

func (r *RoleRouteJsonStorage) DelRoleRoute(role, route string) {
	routes := r.data[role]
	for i, rute := range routes {
		if rute == route {
			r.data[role] = append(routes[:i], routes[i+1:]...)
			delete(r.mp[role], route)
			return
		}
	}
}

func (r *RoleRouteJsonStorage) MatchRoleRoute(role, route string) bool {
	routes, ok := r.mp[role]
	if ok {
		_, ok = routes[route]
		return ok
	} else {
		return false
	}
}
