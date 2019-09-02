// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package env

var Env = &env{}

/**

# 授权加密 key
auth_key: "change this pass"

# 授权过期时间/小时
auth_expire_hour: 168

# casbin 权限模型文件
auth_model_file: "configs/rbac_model.conf"

# 初始管理员用户名
root_user: "root"

# 初始管理员密码
root_password: "root654321"
**/
type env struct {
	AuthKey string `yaml:"auth_key"`
	AuthExpireHour int `yaml:"auth_expire_hour"`
	AuthModelFile string `yaml:"auth_model_file"`

	RootUser string `yaml:"root_user"`
	RootPassword string `yaml:"root_password"`
}