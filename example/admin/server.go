// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
	"github.com/orivil/morgine/bundles/admin"
	admin_model "github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/bundles/utils/api"
	"github.com/orivil/morgine/cfg"
	"github.com/orivil/morgine/x_init"
	"github.com/orivil/morgine/xx"
	"net/http"
)

var env =
`# 授权加密 key
auth_key: "change this pass"

# 授权过期时间/小时
auth_expire_hour: 168

# casbin 权限模型文件
auth_model_file: "configs/rbac_model.conf"

# 开启日志
db_log: true

# mysql postgres
db_dialect: "postgres"

# 数据库地址, 线上项目应该从OS环境变量中获取
db_host: "localhost"

# 数据库监听端口, 线上项目应该从OS环境变量中获取
db_port: ""

# 用户名, 线上项目应该从OS环境变量中获取
db_user: ""

# 密码, 线上项目应该从OS环境变量中获取1
db_password: ""

# 数据库名
db_name: ""

# 表前缀
db_sql_table_prefix: ""

# 最大空闲连接, 支持热重载
db_max_idle_connects: 5

# 最大活动连接, 支持热重载
db_max_opened_connects: 10

# http 服务监听地址
http_addr: ":9090"
`

func main() {
	configs, err := cfg.UnmarshalMap("env.yml", env)
	if err != nil {
		panic(err)
	}
	xx.Use(xx.Cors)
	xx.Handle(http.MethodOptions, "/", nil, func(ctx *xx.Context) {})
	xx.Handle("GET", "/api-data", &xx.Doc{
		Title: "API DATA",
	}, func(ctx *xx.Context) {
		err := ctx.SendJSON(xx.MAP{
			"doc": xx.DefaultServeMux.ApiDoc(),
			"models": api.Models,
		})
		if err != nil {
			ctx.Error(err)
		}
	})
	x_init.Register(configs, admin.Bundle)

	as, err := admin_model.GetRoleAdmins(1, 10, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(as))
	xx.Run(configs.GetStr("http_addr"))
}