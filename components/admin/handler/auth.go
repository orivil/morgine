// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package handler

import "github.com/orivil/morgine/xx"

func Login(method, route string, cdt *xx.Condition) {
	type params struct {
		Username string `param:"username" desc:"用户名"`
		Password string `param:"password" desc:"密码"`
	}
	doc := &xx.Doc{
		Title:     "登录",
		Desc:      "",
		Params:    xx.Params{
			{Type:xx.Form, Schema:&params{}},
		},
		Responses: xx.Responses{
			{
				Description: "将 token 保存起来, 在全局请求头 Header 中都加入 Authorization: token 值",
				Body:xx.JsonData(xx.StatusSuccess, "authorization token"),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		ps := &params{}
		err := ctx.Unmarshal(ps)
		if err != nil {
			xx.HandleError(ctx, err)
		} else {

		}
	})
}
