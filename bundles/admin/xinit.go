// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin

import (
	"github.com/orivil/morgine/bundles/admin/env"
	"github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/bundles/utils/sql"
	"github.com/orivil/morgine/cfg"
)

var Bundle bundle

type bundle struct {
	configs cfg.Configs
}

func (b bundle) Init(configs cfg.Configs) {
	{
		// 链接数据库
		env := &sql.Env{}
		err := configs.Unmarshal(env)
		if err != nil {
			panic(err)
		}
		admin_model.DB, err = env.Connect("admin_")
		if err != nil {
			panic(err)
		}
	}
	{
		// 初始化配置数据
		err := env.Init(configs)
		if err != nil {
			panic(err)
		}
	}
}

func (b bundle) AddRoute() {
	// 注册路由
	registerRoutes()
}

func (b bundle) Run() {
	// 迁移数据库模型
	admin_model.DB.AutoMigrate(&admin_model.Account{})

	// 创建初始账户
	total, err := admin_model.CountAdmins()
	if err != nil {
		panic(err)
	}
	if total == 0 {
		err := admin_model.CreateAdmin(env.Env.RootUser, env.Env.RootPassword)
		if err != nil {
			panic(err)
		}
	}
}
