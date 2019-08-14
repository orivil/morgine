// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package actions

import (
	admin_model "github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/xx"
	"github.com/orivil/morgine/xx/middleware"
)

var Login xx.Action

var Create xx.Action = func(method, route string, controller *xx.Condition) {
	type param struct {
		Username string `required:"用户名不能为空" len:"6-12" len-msg:"用户名必须在6-12个字符之间"`
		Password string `required:"密码不能为空" len:"6-12" len-msg:"密码必须在6-12个字符之间" reg:"\\w|\\_" reg-msg:"密码只能是字母数字或下划线"`
	}
	doc := &xx.Doc {
		Title: "创建管理员",
		Params:xx.Params{
			{
				Type:xx.Form,
				Schema:&param{},
			},
		},
		Responses:xx.Responses{
			{
				Body:xx.MAP{"admin": admin_model.Account{}},
			},
		},
	}
	controller.Handle(method, route, doc, func(ctx *xx.Context) {
		middleware.GetUserIDFromContext(ctx)
	})
}