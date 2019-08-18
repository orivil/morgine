// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package actions

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/orivil/morgine/bundles/admin/env"
	admin_middleware "github.com/orivil/morgine/bundles/admin/middleware"
	admin_model "github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/xx"
	"strconv"
	"time"
)

var Login xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		Username string `required:"用户名不能为空"`
		Password string `required:"密码不能为空"`
	}
	doc := &xx.Doc{
		Title: "获得登录授权",
		Params: xx.Params{
			{
				Type:   xx.Form,
				Schema: &params{},
			},
		},
		Responses: xx.Responses{
			xx.MessageResponse(xx.MsgTypeWarning),
			{
				Body: xx.MAP{"authorization": "token string"},
			},
		},
	}
	expire := time.Duration(env.Env.AuthExpireHour) * time.Hour
	controller.Handle(method, route, doc, func(ctx *xx.Context) {
		p := &params{}
		err := ctx.Unmarshal(p)
		if err != nil {
			xx.HandleUnmarshalError(err, ctx)
		} else {
			id, err := admin_model.SignIn(p.Username, p.Password)
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
					ExpiresAt: time.Now().Add(expire).Unix(),
					Id:        strconv.Itoa(id),
				})
				auth, err := token.SignedString([]byte(env.Env.AuthKey))
				if err != nil {
					ctx.Error(err)
				} else {
					ctx.SendJSON(xx.MAP{"authorization": auth})
				}
			}
		}
	})
}

var ChangePassword xx.Action = func(method, route string, controller *xx.Condition) {
	type param struct {
		Username string `required:"用户名不能为空"`
		OldPassword string `required:"旧密码不能为空"`
		// TODO:
		NewPassword string `required:"新密码不能为空" len:"6-12" len-msg:"密码必须在6-12个字符之间" reg:"^[\\w|\\_]+$" reg-msg:"密码只能是字母数字或下划线"`
	}
	doc := &xx.Doc {
		Title: "更改管理员密码",
		Params:xx.Params {
			{
				Type:xx.Form,
				Schema:&param{},
			},
		},
	}
	controller.Handle(method, route, doc, func(ctx *xx.Context) {
		p := &param{}
		err := ctx.Unmarshal(p)
		if err != nil {
			xx.HandleUnmarshalError(err, ctx)
		} else {
			if id, ok := admin_middleware.GetUserIDFromContext(ctx); ok {
				err := admin_model.UpdatePassword(id, p.Username, p.OldPassword, p.NewPassword)
				if err != nil {
					ctx.MsgWarning(err.Error())
				} else {
					ctx.MsgSuccess("修改成功")
				}
			} else {
				ctx.MsgWarning("用户未登录")
			}
		}
	})
}