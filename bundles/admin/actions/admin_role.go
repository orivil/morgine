// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package actions

import (
	admin_model "github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/xx"
)

var CountAdminRoles xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		AdminID int `required:""`
	}
	doc := xx.Doc {
		Title: "统计管理员角色数量",
		Params:xx.Params {
			{
				Type: xx.Query,
				Schema: &params{},
			},
		},
		Responses:xx.Responses {
			{
				Body: xx.MAP{"total": 5},
			},
		},
	}
	controller.Handle(method, route, &doc, func(ctx *xx.Context) {
		p := &params{}
		err := ctx.Unmarshal(p)
		if err != nil {
			xx.HandleUnmarshalError(err, ctx)
		} else {
			roles, err := admin_model.CountAdminRoles(p.AdminID)
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.SendJSON(xx.MAP{"roles": roles})
			}
		}
	})
}

var GetAdminRoles xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		AdminID int `required:""`
		Limit int `required:"" num:"1-200"`
		Offset int
	}
	doc := xx.Doc {
		Title: "获得管理员角色列表",
		Params:xx.Params {
			{
				Type: xx.Query,
				Schema: &params{},
			},
		},
		Responses:xx.Responses {
			{
				Body: xx.MAP{"roles": []*admin_model.AdminRole{{}}},
			},
		},
	}
	controller.Handle(method, route, &doc, func(ctx *xx.Context) {
		p := &params{}
		err := ctx.Unmarshal(p)
		if err != nil {
			xx.HandleUnmarshalError(err, ctx)
		} else {
			roles, err := admin_model.GetAdminRoles(p.AdminID, p.Limit, p.Offset)
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.SendJSON(xx.MAP{"roles": roles})
			}
		}
	})
}

var SetAdminRole xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		AdminID int `required:""`
		RoleID int `required:""`
	}
	doc := xx.Doc {
		Title: "设置管理员角色",
		Desc: "如果角色已存在会报错",
		Params:xx.Params {
			{
				Type: xx.Query,
				Schema: &params{},
			},
		},
		Responses:xx.Responses {
			xx.MessageResponse(xx.MsgTypeSuccess),
		},
	}
	controller.Handle(method, route, &doc, func(ctx *xx.Context) {
		p := &params{}
		err := ctx.Unmarshal(p)
		if err != nil {
			xx.HandleUnmarshalError(err, ctx)
		} else {
			err := admin_model.SetAdminRole(p.AdminID, p.RoleID)
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.MsgSuccess("已保存")
			}
		}
	})
}

var DelAdminRoles xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		AdminID int `required:"" desc:"管理员ID"`
		RoleIDs []int `required:"" desc:"角色ID"`
	}
	doc := xx.Doc {
		Title: "移除管理员角色",
		Params:xx.Params {
			{
				Type: xx.Query,
				Schema: &params{},
			},
		},
		Responses:xx.Responses {
			xx.MessageResponse(xx.MsgTypeSuccess),
		},
	}
	controller.Handle(method, route, &doc, func(ctx *xx.Context) {
		p := &params{}
		err := ctx.Unmarshal(p)
		if err != nil {
			xx.HandleUnmarshalError(err, ctx)
		} else {
			err := admin_model.RemoveAdminRoles(p.AdminID, p.RoleIDs)
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.MsgSuccess("已移除")
			}
		}
	})
}