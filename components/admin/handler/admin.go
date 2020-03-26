// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package handler

import (
	"github.com/orivil/morgine/components/admin/auth"
	"github.com/orivil/morgine/components/admin/models"
	"github.com/orivil/morgine/components/admin/models/api"
	"github.com/orivil/morgine/utils/sql"
	"github.com/orivil/morgine/xx"
)

func Login(method, route string, cdt *xx.Condition) {
	type params struct {
		Username string `param:"username" desc:"用户名"`
		Password string `param:"password" desc:"密码"`
	}
	doc := &xx.Doc {
		Title:     "登录",
		Desc:      "",
		Params:    xx.Params {
			{Type:xx.Form, Schema:&params{}},
		},
		Responses: xx.Responses {
			{
				Description: "将 token 保存起来, 在全局请求头 Header 中都加入 Authorization: token 值",
				Body:xx.JsonData(xx.StatusSuccess, xx.MAP{
					"user": models.Admin{},
					"token": "authorization token",
				}),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		ps := &params{}
		err := ctx.Unmarshal(ps)
		if err != nil {
			xx.HandleError(ctx, err)
		} else {
			var admin *models.Admin
			admin, err = api.LoginAdmin(ps.Username, ps.Password)
			if err != nil {
				if err == api.ErrLoginUsernameIncorrect {
					xx.SendMessage(ctx, xx.MsgTypeError, "用户名错误")
				} else if err == api.ErrMismatchedPassword {
					xx.SendMessage(ctx, xx.MsgTypeError, "密码错误")
				} else {
					xx.HandleError(ctx, err)
				}
			} else {
				var token []byte
				token, err = auth.EncryptToken(admin.ID)
				if err != nil {
					xx.HandleError(ctx, err)
				} else {
					xx.SendJson(ctx, xx.StatusSuccess, xx.MAP {
						"user": admin,
						"token": string(token),
					})
				}
			}
		}
	})
}

func GetUserInfo(method, route string, cdt *xx.Condition) {
	doc := &xx.Doc {
		Title:     "获得登录用户信息",
		Desc:      "",
		Params:    nil,
		Responses: xx.Responses{
			{
				Description: "登录成功",
				Body: xx.JsonData(xx.StatusSuccess, &models.Admin{}),
			},
			{
				Description: "登录状态已失效",
				Body: xx.JsonData(xx.StatusNotFound, nil),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		id := auth.GetAdminID(ctx)
		admin := api.GetAdminByID(id)
		if admin != nil {
			xx.SendJson(ctx, xx.StatusSuccess, admin)
		} else {
			xx.SendJson(ctx, xx.StatusNotFound, nil)
		}
	})
}

func CreateSubAccount(method, route string, cdt *xx.Condition) {
	type params struct {
		ParentID int `desc:"子账号的父ID，如果不提供该参数则创建当前登录账号的子账号"`
		Username   string      `reg:"^[a-zA-Z0-9]{4,16}$" desc:"账号，4-16字母或数字"`
		Nickname   string      `desc:"昵称"`
		Password   string      `reg:"^[a-zA-Z0-9]{8,16}$" desc:"密码，8-16字母或数字"`
		Super      sql.Boolean `desc:"1-超级管理员 2-普通管理员"`
	}
	doc := &xx.Doc{
		Title:     "创建子账号或子孙账号",
		Desc:      "当前登录账号必须是超级管理员权限",
		Params:    xx.Params{
			{
				Type:xx.Form,
				Schema:&params{},
			},
		},
		Responses: xx.Responses{
			{
				Description: "返回子账号信息",
				Body:xx.JsonData(xx.StatusSuccess, &models.Admin{}),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		loginID := auth.GetAdminID(ctx)
		if loginID > 0 {
			ps := &params{}
			err := ctx.Unmarshal(ps)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				admin := &models.Admin{
					Username:   ps.Username,
					Nickname:   ps.Nickname,
					Password:   ps.Password,
					Super:      ps.Super,
				}
				err = api.CreateSubAdmin(loginID, ps.ParentID, admin)
				if err != nil {
					xx.HandleError(ctx, err)
				} else {
					xx.SendJson(ctx, xx.StatusSuccess, admin)
				}
			}
		}
	})
}

func DelSubAccount(method, route string, cdt *xx.Condition) {
	type params struct {
		ID int `desc:"子账号ID"`
	}
	doc := &xx.Doc{
		Title:     "删除子账号",
		Desc:      "只能删除当前登录账号的子账号",
		Params:    xx.Params{
			{
				Type: xx.Query,
				Schema: &params{},
			},
		},
		Responses: xx.Responses {
			{
				Description: "删除成功",
				Body: xx.JsonData(xx.StatusSuccess, nil),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		parentID := auth.GetAdminID(ctx)
		if parentID > 0 {
			ps := &params{}
			err := ctx.Unmarshal(ps)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				err = api.DelSubAccount(parentID, ps.ID)
				if err != nil {
					xx.SendMessage(ctx, xx.MsgTypeError, err.Error())
				} else {
					xx.SendJson(ctx, xx.StatusSuccess, nil)
				}
			}
		}
	})
}

func UpdatePassword(method, route string, cdt *xx.Condition)  {
	type params struct {
		SubID int `param:"sub_id" desc:"子账号ID，如果不提供该参数，则修改当前登录用户的密码，否则修改子账号密码"`
		Password string `param:"password" reg:"^[a-zA-Z0-9]{8,16}$" desc:"密码，8-16字母或数字"`
	}
	doc := &xx.Doc{
		Title:     "更新当前账号密码",
		Desc:      "",
		Params:    xx.Params{
			{
				Type:xx.Form,
				Schema:&params{},
			},
		},
		Responses: xx.Responses{
			{
				Body:xx.JsonData(xx.StatusSuccess, nil),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		id := auth.GetAdminID(ctx)
		if id > 0 {
			ps := &params{}
			err := ctx.Unmarshal(ps)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				err = api.UpdateAdminPassword(id, ps.SubID, ps.Password)
				if err != nil {
					xx.HandleError(ctx, err)
				} else {
					xx.SendJson(ctx, xx.StatusSuccess, nil)
				}
			}
		}
	})
}

func UpdateAdminInfo(method, route string, cdt *xx.Condition) {
	type params struct {
		SubID int `desc:"子账号ID，如果不提供该参数则修改当前登录账号的信息"`
		Nickname string `desc:"昵称"`
	}
	doc := &xx.Doc{
		Title:     "更新管理员信息",
		Desc:      "",
		Params:    xx.Params{
			{Type:xx.Form, Schema:&params{}},
		},
		Responses: nil,
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		parentID := auth.GetAdminID(ctx)
		if parentID > 0 {
			ps := &params{}
			err := ctx.Unmarshal(ps)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				err = api.UpdateAdminInfo(parentID, ps.SubID, &models.Admin{Nickname: ps.Nickname})
				if err != nil {
					xx.HandleError(ctx, err)
				} else {
					xx.SendJson(ctx, xx.StatusSuccess, nil)
				}
			}
		}
	})
}

func GetAllSubAccounts(method, route string, cdt *xx.Condition) {
	doc := &xx.Doc{
		Title:     "获取所有子账号",
		Desc:      "",
		Params:    nil,
		Responses: xx.Responses {
			{
				Body: xx.JsonData(xx.StatusSuccess, []*api.Account{{}}),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		id := auth.GetAdminID(ctx)
		accounts := api.GetSubAdmins(id)
		xx.SendJson(ctx, xx.StatusSuccess, accounts)
	})
}