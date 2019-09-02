// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin_middleware

import (
	"github.com/casbin/casbin"
	gormadapter "github.com/casbin/gorm-adapter"
	"github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/orivil/morgine/bundles/admin/env"
	admin_model "github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/xx"
	"strconv"
	"time"
)


var Enforcer *casbin.Enforcer

func InitEnforcer(modelFile string, db *gorm.DB) error {
	a, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return err
	}
	e, err := casbin.NewEnforcer(modelFile, a)
	if err != nil {
		return err
	}
	err = e.LoadPolicy()
	if err != nil {
		return err
	}
	Enforcer = e
	return nil
}

var userIDContextKey = "user-id"

var userRolesContextKey = "user-roles"

var Auth = func() *xx.Handler {
	type ps struct {
		Authorization string `required:"未提供授权码" desc:"该数据需通过登录接口获得"`
	}
	doc := &xx.Doc {
		Title: "User Auth",
		Desc: "用户登录中间件",
		Params: xx.Params {
			{
				Type:   xx.Header,
				Schema: &ps{},
			},
		},
		Responses: xx.Responses {
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
				token, err := jwt.ParseWithClaims(p.Authorization, &adminClaims{}, func(token *jwt.Token) (i interface{}, e error) {
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
					claims := token.Claims.(*adminClaims)
					ok, err := Enforcer.Enforce(claims.Admin.Username)
					if err != nil {
						ctx.Error(err)
					} else {
						if ok {
							ctx.Set(userIDContextKey, claims.Admin)
						} else {

						}
					}
					ctx.Set(userRolesContextKey, claims.Roles)
				}
			}
		},
	}
}()

type adminClaims struct {
	Admin *admin_model.Admin
	//Roles []int
	jwt.StandardClaims
}

// 从上下文(Context)中获取管理员ID，需要在授权中间件处理之后调用
func GetUserIDFromContext(ctx *xx.Context) (int, bool) {
	id, ok := ctx.Get(userIDContextKey).(int)
	return id, ok
}

// 从上下文(Context)中获取管理员角色列表, 需要在授权中间件处理之后调用
func GetUserRolesFromContext(ctx *xx.Context) ([]int, bool) {
	ids, ok := ctx.Get(userRolesContextKey).([]int)
	return ids, ok
}

func NewToken(admin *admin_model.Admin) (token string, err error) {
	jt := jwt.NewWithClaims(jwt.SigningMethodHS256, adminClaims{
		Admin:admin,
		StandardClaims: jwt.StandardClaims {
			ExpiresAt: time.Now().Add(time.Duration(env.Env.AuthExpireHour) * time.Hour).Unix(),
			Id:        strconv.Itoa(admin.ID),
		},
	})
	return jt.SignedString([]byte(env.Env.AuthKey))
}