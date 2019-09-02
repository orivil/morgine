// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin

import (
	"github.com/orivil/morgine/bundles/admin/env"
	admin_middleware "github.com/orivil/morgine/bundles/admin/middleware"
	"github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/bundles/utils/api"
	"github.com/orivil/morgine/bundles/utils/sql"
	"github.com/orivil/morgine/cfg"
)

var Bundle bundle

type bundle struct {
	configs cfg.Configs
}

func (b bundle) Init(configs cfg.Configs) {
	{
		// 初始化配置数据
		err := configs.Unmarshal(env.Env)
		if err != nil {
			panic(err)
		}
	}
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
		// 迁移数据库模型
		admin_model.DB = api.AutoMigrate(admin_model.DB, &admin_model.Admin{}, "管理员账号数据模型")
		admin_model.DB = api.AutoMigrate(admin_model.DB, &admin_model.Role{}, "角色数据模型")
		admin_model.DB = api.AutoMigrate(admin_model.DB, &admin_model.AdminRole{}, "管理员-角色对照表")

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

	{
		// 初始化角色权限控制器
		err := admin_middleware.InitEnforcer(env.Env.AuthModelFile, admin_model.DB)
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

}