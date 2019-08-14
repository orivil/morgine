// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin

import (
	"github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/bundles/utils/sql"
	"github.com/orivil/morgine/cfg"
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
		admin_model.DB, err = env.Connect("admin_")
		if err != nil {
			panic(err)
		}
	}
	{
		// 初始化 middleware 及 action
		env := &Env{}
		err := configs.Unmarshal(env)
		if err != nil {
			panic(err)
		}
		err = env.Init()
		if err != nil {
			panic(err)
		}
	}
}

func (b bundle) AddRoute() {
	registerRoutes()
}

func (b bundle) Run() {
	admin_model.DB.AutoMigrate(&admin_model.Account{})
}
