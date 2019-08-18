// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package env

import "github.com/orivil/morgine/cfg"

var Env = &env{}

/**

# 授权加密 key
auth_key: "change this pass"

# 授权过期时间/小时
auth_expire_hour: 168

# 初始管理员用户名
root_user: "root"

# 初始管理员密码
root_password: "root654321"
**/
type env struct {
	AuthKey string `yaml:"auth_key"`
	AuthExpireHour int `yaml:"auth_expire_hour"`

	RootUser string `yaml:"root_user"`
	RootPassword string `yaml:"root_password"`
}

func Init(configs cfg.Configs) error {
	return configs.Unmarshal(Env)
}