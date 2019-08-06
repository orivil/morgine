// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"github.com/orivil/morgine/xx"
	"net/http"
)

var (
	sitesService       = xx.NewTagName("站点服务")
	accountsService    = xx.NewTagName("管理员账号服务")
	filesService       = xx.NewTagName("文件服务")
	imageLabelsService = xx.NewTagName("图片标签服务")
	imagesService      = xx.NewTagName("图片服务")
)

var tags = xx.ApiTags{
	{
		Name: sitesService,
		Subs: xx.ApiTags{
			{
				Name: accountsService,
				Subs: xx.ApiTags{
					{Name: filesService},
					{Name: imageLabelsService},
					{Name: imagesService},
				},
			},
		},
	},
}

func main() {
	xx.Use(xx.Cors)
	xx.Handle(http.MethodOptions, "/", func(ctx *xx.Context) {}, nil)
	xx.Handle("GET", "/foo", func(ctx *xx.Context) {
		ctx.WriteString("bar")
	}, &xx.Doc{
		Title: "FOO BAR",
	})
	xx.Handle("GET", "/{mp}.txt", func(ctx *xx.Context) {
		ctx.WriteString(ctx.Path().Get("mp"))
	}, &xx.Doc{
		Title: "MP text",
	})
	xx.Handle("GET", "/api-data", func(ctx *xx.Context) {
		err := ctx.SendJSON(xx.MAP{"doc": xx.DefaultServeMux.ApiDoc()})
		if err != nil {
			panic(err)
		}
	}, &xx.Doc{
		Title: "API DATA",
	})
	group := xx.NewGroup(tags)
	group = group.Use(mustLogin)
	accountController := group.Controller(accountsService)
	handleLogin("GET", "/login", accountController)
	xx.Run()
}

func handleLogin(method, route string, group *xx.RouteGroup) {
	type param struct {
		Username string
		Password string
	}
	doc := &xx.Doc{
		Title: "登录管理员",
		Params: xx.Params{
			{
				Type:   xx.Query,
				Schema: &param{},
			},
		},
	}
	group.Handle(method, route, func(ctx *xx.Context) {
		p := &param{}
		err := ctx.Unmarshal(p)
		if err != nil {
			ctx.MsgWarning(err.Error())
		} else {
			ctx.SendJSON(p)
		}
	}, doc)
}

var mustLogin = func() *xx.Handler {
	type authorization struct {
		Authorization string
	}
	return &xx.Handler{
		Doc: &xx.Doc{
			Title: "Admin Must Login",
			Desc:  "管理员登录权限",
			Params: xx.Params{
				{
					Type:   xx.Header,
					Schema: &authorization{},
				},
			},
			Responses: xx.Responses{{
				Code: http.StatusForbidden,
				Body: xx.MsgData(xx.MsgTypeWarning, "需要管理员权限"),
			}},
		},
		HandleFunc: func(ctx *xx.Context) {
			p := &authorization{}
			ctx.Unmarshal(p)
			if p.Authorization == "" {
				ctx.MsgWarning("需要管理员权限")
			}
		},
	}
}()
