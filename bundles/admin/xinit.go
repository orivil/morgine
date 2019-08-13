// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin

import (
	"github.com/orivil/morgine/bundles/admin/middleware"
	"github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/bundles/utils/sql"
	"github.com/orivil/morgine/cfg"
	"github.com/orivil/morgine/xx"
	middleware2 "github.com/orivil/morgine/xx/middleware"
)

var Bundle bundle

type bundle int

func (b bundle) Init(configs cfg.Configs) {
	{
		// 链接数据库
		env := &sql.Env{}
		err := configs.Unmarshal(env)
		if err != nil {
			panic(err)
		}
		model.DB, err = env.Connect("admin_")
		if err != nil {
			panic(err)
		}
	}
	{
		// 初始化中间件
		authKey := configs.GetStr("auth_key")
		if authKey == "" {
			panic("auth_key is empty")
		}
		middleware.AdminAuth = middleware2.NewJWT([]byte(authKey))
	}
}

func (b bundle) AddRoute() {
	var (
		service
	)
	var tags = xx.ApiTags{

	}
}

func (b bundle) Run() {
	model.DB.AutoMigrate(&model.Admin{})
}
