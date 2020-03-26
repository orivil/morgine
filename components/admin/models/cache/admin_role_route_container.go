// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cache

import (
	"github.com/orivil/morgine/utils/cache"
	"strconv"
	"time"
)

type RouteStorage interface {
	IsAdminAllowedRoute(adminID int, route string) (ok bool)
}

type AdminRoleRoutesContainer struct {
	container *cache.Container
	storage RouteStorage
}

func NewAdminRoleRoutes(storage RouteStorage) *AdminRoleRoutesContainer {
	return &AdminRoleRoutesContainer{
		container: cache.NewContainer(),
		storage: storage,
	}
}

func (r *AdminRoleRoutesContainer) Check(adminID int, route string) bool {
	key := strconv.Itoa(adminID) + route
	vue, ok := r.container.Get(key).(bool)
	if !ok {
		vue = r.storage.IsAdminAllowedRoute(adminID, route)
		expireAt := time.Now().Add(5 * time.Minute)
		r.container.Set(key, vue, &expireAt)
	}
	return vue
}