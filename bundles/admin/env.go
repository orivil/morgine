// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin

import (
	"errors"
	"github.com/orivil/morgine/bundles/admin/actions"
	admin_middleware "github.com/orivil/morgine/bundles/admin/middleware"
	admin_model "github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/xx/middleware"
	"time"
)

// 配置示例
var env =
`
# 授权加密 key
auth_key: "change this pass"

# 授权过期时间/小时
auth_expire_hour: 168
`

type Env struct {
	AuthKey string `yaml:"auth_key"`
	AuthExpireHour int `yaml:"auth_expire_hour"`
}

func (e *Env) Init() error {
	if len(e.AuthKey) == 0 {
		return errors.New("auth_key is empty")
	}
	// 初始化 middleware
	admin_middleware.Auth = middleware.NewJWT([]byte(e.AuthKey))

	// 初始化登录 action
	actions.Login = middleware.NewLoginHandler([]byte(e.AuthKey), time.Duration(e.AuthExpireHour) * time.Hour, admin_model.SignIn)
	return nil
}