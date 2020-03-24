// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package cache

import (
	"github.com/orivil/morgine/utils/cache"
	"strconv"
)

type RoleRouteBackend struct {
	container *cache.Container
}

func (r *RoleRouteBackend) Check(adminID int, route string) bool {
	key := strconv.Itoa(adminID) + route
	vue := r.container.Get(key)
	if vue == nil {

	}
}