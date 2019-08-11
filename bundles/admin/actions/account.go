// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package actions

import (
	"encoding/base64"
	"github.com/orivil/morgine/bundles/utils/crypto"
	"github.com/orivil/morgine/xx"
)

func getAuthoration(method, route string, r *xx.RouteGroup) {
	type ps struct {
		Username string `required:"用户名不能为空"`
		Password string `required:"密码不能为空"`
	}
	doc := &xx.Doc{
		Title: "获得授权码",
		Params: xx.Params{
			{
				Type:   xx.Query,
				Schema: &ps{},
			},
		},
		Responses: xx.Responses{
			{
				Body: xx.MAP{"authorization": "Basic xxxxxxxxxxx"},
			},
		},
	}
	r.Handle(method, route, doc, func(ctx *xx.Context) {
		p := &ps{}
		err := ctx.Unmarshal(p)
		if err != nil {
			ctx.MsgWarning(err.Error())
		} else {
			username, err := crypto.EncryptString(p.Username)
			if err != nil {
				ctx.Message(xx.MsgTypeWarning, err.Error())
			} else {
				password, err := crypto.EncryptString(p.Password)
				if err != nil {
					ctx.Message(xx.MsgTypeWarning, err.Error())
				} else {
					auth := username + ":" + password
					auth = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
					ctx.SendJSON(xx.MAP{"authorization": auth})
				}
			}
		}
	})
}
