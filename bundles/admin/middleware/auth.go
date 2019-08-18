// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin_middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/orivil/morgine/bundles/admin/env"
	"github.com/orivil/morgine/xx"
	"strconv"
)

var userIDContextKey = "user-id"

var Auth = func() *xx.Handler {
	type ps struct {
		Authorization string `required:"未提供授权码" desc:"该数据需通过登录接口获得"`
	}
	doc := &xx.Doc{
		Title: "用户登录中间件",
		Params: xx.Params{
			{
				Type:   xx.Header,
				Schema: &ps{},
			},
		},
		Responses: xx.Responses{
			xx.MessageResponse(xx.MsgTypeWarning),
		},
	}
	return &xx.Handler {
		Doc: doc,
		HandleFunc: func(ctx *xx.Context) {
			p := &ps{}
			err := ctx.Unmarshal(p)
			if err != nil {
				xx.HandleUnmarshalError(err, ctx)
			} else {
				token, err := jwt.ParseWithClaims(p.Authorization, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
					return []byte(env.Env.AuthKey), nil
				})
				if err != nil {
					if ve, ok := err.(*jwt.ValidationError); ok {
						if ve.Errors&jwt.ValidationErrorExpired != 0 {
							ctx.MsgWarning("管理员授权过期")
							return
						}
					}
					ctx.MsgWarning(err.Error())
				} else {
					claims := token.Claims.(*jwt.StandardClaims)
					id, err := strconv.Atoi(claims.Id)
					if err != nil {
						ctx.Error(err)
					} else {
						ctx.Set(userIDContextKey, id)
					}
				}
			}
		},
	}
}()

// 从上下文(Context)中获取管理员ID，需要在授权中间件处理之后调用
func GetUserIDFromContext(ctx *xx.Context) (int, bool) {
	id, ok := ctx.Get(userIDContextKey).(int)
	return id, ok
}