// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/orivil/morgine/xx"
	"strconv"
	"strings"
	"time"
)

const AuthScheme = "Bearer "

// 新建授权中间件, key 参数用于数据加密, 如果其他中间件也需要
// 获取登录用户信息, 则需要设置在授权中间件之后
func NewJWT(key []byte) *xx.Handler {
	type ps struct {
		Authorization string `desc:"该数据需通过登录接口获得"`
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
	return &xx.Handler{
		Doc: doc,
		HandleFunc: func(ctx *xx.Context) {
			p := &ps{}
			err := ctx.Unmarshal(p)
			if err != nil {
				ctx.Error(err)
			} else {
				if p.Authorization == "" {
					ctx.MsgWarning("管理员未登录")
				} else {
					auth := strings.TrimPrefix(p.Authorization, AuthScheme)
					token, err := jwt.ParseWithClaims(auth, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
						return key, nil
					})
					if err != nil {
						ctx.Error(err)
					} else {
						if !token.Valid {
							ctx.MsgWarning("管理员授权失败")
						} else {
							claims := token.Claims.(*jwt.StandardClaims)
							if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
								ctx.MsgWarning("管理员授权过期")
							} else {
								id, err := strconv.Atoi(claims.Id)
								if err != nil {
									ctx.Error(err)
								} else {
									ctx.Set(userIDContextKey, id)
								}
							}
						}
					}
				}
			}
		},
	}
}

var userIDContextKey = "user-id"

// 从上下文(Context)中获取管理员ID，需要在授权中间件处理之后调用
func GetUserIDFromContext(ctx *xx.Context) (int, bool) {
	id, ok := ctx.Get(userIDContextKey).(int)
	return id, ok
}

// UserIDProvider 用于验证账户密码, 验证成功则返回用户 id, 验证失败返回 err
type UserIDProvider func(username, password string) (id int, err error)

func NewLoginHandler(key []byte, expire time.Duration, provider UserIDProvider) xx.Action {
	return func(method, route string, rg *xx.RouteGroup) {
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
					Body: xx.MAP{"authorization": AuthScheme + "xxxxxxxxxxxxx"},
				},
			},
		}
		rg.Handle(method, route, doc, func(ctx *xx.Context) {
			p := &params{}
			err := ctx.Unmarshal(p)
			if err != nil {
				xx.HandleUnmarshalError(err, ctx)
			} else {
				id, err := provider(p.Username, p.Password)
				if err != nil {
					ctx.MsgWarning(err.Error())
				} else {
					token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
						ExpiresAt: time.Now().Add(expire).Unix(),
						Id:        strconv.Itoa(id),
					})
					auth, err := token.SignedString(key)
					if err != nil {
						ctx.Error(err)
					} else {
						ctx.SendJSON(xx.MAP{"authorization": AuthScheme + auth})
					}
				}
			}
		})
	}
}
