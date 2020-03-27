// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

// +build ignore

package main

import (
	"github.com/orivil/morgine/param"
	"github.com/orivil/morgine/xx"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var (
	sitesService       = xx.NewTagName("站点服务")
	adminService       = xx.NewTagName("管理员服务")
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
				Name: adminService,
				Subs: xx.ApiTags{
					{Name: accountsService},
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
	xx.Handle(http.MethodOptions, "/", nil, func(ctx *xx.Context) {})
	xx.Handle("GET", "/foo", &xx.Doc {
		Title: "FOO BAR",
	}, func(ctx *xx.Context) {
		ctx.WriteString("bar")
	})
	{
		type mp struct {
			mp string
		}
		xx.Handle("GET", "/{mp}.txt", &xx.Doc{
			Title: "MP text",
			Params: xx.Params{
				{
					Type:   xx.Path,
					Schema: &mp{},
				},
			},
		}, func(ctx *xx.Context) {
			p := &mp{}
			err := ctx.Unmarshal(p)
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.WriteString(p.mp)
			}
		})
	}
	xx.Handle("GET", "/api-data", &xx.Doc{
		Title: "API DATA",
	}, func(ctx *xx.Context) {
		ctx.SendJSON(xx.MAP{"doc": xx.DefaultServeMux.ApiDoc()})
	})
	group := xx.NewGroup(tags)
	//group = group.Use(mustLogin)
	accountController := group.Controller(accountsService)
	handleLogin("POST", "/login", accountController)
	xx.Run(":9090")
}

func handleLogin(method, route string, group *xx.Condition) {
	type pm struct {
		Username     string `required:""`
		Password     string `desc:"密码"`
		String       string
		Int          int
		Int32        int32
		Int64        int64
		Float32      float32
		Float64      float64
		Bool         bool
		File         param.FileHandler
		Time         *time.Time
		SliceString  []string
		SliceInt     []int
		SliceInt32   []int32
		SliceInt64   []int64
		SliceFloat32 []float32
		SliceFloat64 []float64
		SliceBool    []bool
	}
	doc := &xx.Doc{
		Title: "登录管理员",
		Params: xx.Params{
			{
				Type: xx.Form,
				Schema: &pm{
					SliceInt: []int{1, 2, 3},
					File: func(field string, header *multipart.FileHeader) error {
						fs, err := header.Open()
						if err != nil {
							return err
						}
						defer fs.Close()
						data, err := ioutil.ReadAll(fs)
						if err != nil {
							return err
						}
						return ioutil.WriteFile(filepath.Join("imgs", header.Filename), data, os.ModePerm)
					},
				},
			},
		},
	}
	group.Handle(method, route, doc, func(ctx *xx.Context) {
		p := &pm{}
		err := ctx.Unmarshal(p)
		if err != nil {
			ctx.MsgWarning(err.Error())
		} else {
			ctx.SendJSON(p)
		}
	})
}

var mustLogin = func() *xx.Handler {
	type params struct {
		Authorization string `required:"用户未登录"`
	}
	return &xx.Handler{
		Doc: &xx.Doc{
			Title: "Admin Must Login",
			Desc:  "管理员登录权限",
			Params: xx.Params{
				{
					Type:   xx.Header,
					Schema: &params{},
				},
			},
			Responses: xx.Responses{
				xx.MessageResponse(xx.MsgTypeWarning),
			},
		},
		HandleFunc: func(ctx *xx.Context) {
			p := &params{}
			err := ctx.Unmarshal(p)
			if err != nil {
				xx.HandleUnmarshalError(err, ctx)
			} else {
				if p.Authorization == "" {
					ctx.MsgWarning("需要管理员权限")
				}
			}
		},
	}
}()
