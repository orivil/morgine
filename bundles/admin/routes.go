// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin

import (
	"github.com/orivil/morgine/bundles/admin/actions"
	"github.com/orivil/morgine/xx"
)

var adminService = xx.NewTagName("管理员服务")

var tags = xx.ApiTags{
	{
		Name: adminService,
	},
}

func registerRoutes() {
	group := xx.NewGroup(tags)
	handleAdmin(group.Controller(adminService))
}

func handleAdmin(g *xx.Condition)  {
	actions.Login("GET", "/login", g)
}
