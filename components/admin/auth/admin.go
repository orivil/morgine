// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package auth

import (
	"github.com/orivil/morgine/components/admin/env"
	"github.com/orivil/morgine/xx"
	"strconv"
	"strings"
	"time"
)

var Admin = func() *xx.Handler {
	type param struct {
		Authorization string `desc:"授权后获得的 Token"`
	}
	return &xx.Handler {
		Doc: &xx.Doc {
			Title:     "Admin Authorization",
			Desc:      "管理员授权过滤器",
			Params:    xx.Params{
				{
					Type:xx.Header,
					Schema: &param{},
				},
			},
			Responses: xx.Responses {
				{
					Description: "未登录",
					Body: xx.JsonData(xx.StatusUnauthorized, nil),
				},
				{
					Description: "Token 格式不正确",
					Body: xx.JsonData(xx.StatusTokenInvalid, nil),
				},
				{
					Description: "Token 已过期",
					Body: xx.JsonData(xx.StatusTokenExpired, nil),
				},
			},
		},
		HandleFunc: func(ctx *xx.Context) {
			ps := &param{}
			err := ctx.Unmarshal(ps)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				if ps.Authorization == "" {
					xx.SendJson(ctx, xx.StatusUnauthorized, nil)
				} else {
					id, code := decryptToken(ps.Authorization)
					if code != xx.StatusSuccess {
						xx.SendJson(ctx, code, nil)
					} else {
						setAdminToken(ps.Authorization, ctx)
						setAdminID(id, ctx)
					}
				}
			}
		},
	}
}()

const sep = "|"

func EncryptToken(adminID int) (token []byte, err error) {
	now := strconv.FormatInt(time.Now().Unix(), 10)
	id := strconv.Itoa(adminID)
	data := []byte(now + sep + id)
	return env.AesCrypto.Encrypt(data)
}

func decryptToken(token string) (id int, code xx.StatusCode) {
	var data []byte
	var err error
	data, err = env.AesCrypto.Decrypt([]byte(token))
	if err != nil {
		return 0, xx.StatusTokenInvalid
	}
	str := string(data)
	if strings.Count(str, sep) != 1 {
		return 0, xx.StatusTokenInvalid
	}
	strs := strings.Split(str, sep)
	var start int64
	start, err = strconv.ParseInt(strs[0], 10, 64)
	if err != nil {
		return 0, xx.StatusTokenInvalid
	}
	if time.Unix(start, 0).AddDate(0, 0, env.Config.AuthTokenExpiredDay).Before(time.Now()) {
		return 0, xx.StatusTokenExpired
	}
	id, err = strconv.Atoi(strs[1])
	if err != nil {
		return 0, xx.StatusTokenInvalid
	} else {
		return id, xx.StatusSuccess
	}
}

var adminIDKey = "_admin_id"
var adminTokenKey = "_admin_token"

func setAdminID(id int, ctx *xx.Context) {
	ctx.Set(adminIDKey, id)
}

func setAdminToken(token string, ctx *xx.Context) {
	ctx.Set(adminTokenKey, token)
}

func GetAdminID(ctx *xx.Context) (id int) {
	id, _ = ctx.Get(adminIDKey).(int)
	return
}

func GetAdminToken(ctx *xx.Context) (token string) {
	token, _ = ctx.Get(adminTokenKey).(string)
	return
}