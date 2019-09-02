// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package actions

import (
	admin_model "github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/xx"
)

var CountRoles xx.Action = func(method, route string, controller *xx.Condition) {
	doc := xx.Doc {
		Title: "统计角色数量",
		Responses:xx.Responses{
			{
				Body: xx.MAP{"total": 3},
			},
		},
	}
	controller.Handle(method, route, &doc, func(ctx *xx.Context) {
		total, err := admin_model.CountRoles()
		if err != nil {
			ctx.MsgWarning(err.Error())
		} else {
			ctx.SendJSON(xx.MAP{"total": total})
		}
	})
}

var GetRoles xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		Limit int `required:"" num:"1-100"`
		Offset int
	}
	doc := xx.Doc {
		Title: "获得角色列表",
		Params:xx.Params {
			{
				Type: xx.Query,
				Schema: &params{},
			},
		},
		Responses:xx.Responses {
			{
				Body: xx.MAP{"roles": []*admin_model.Role{{}}},
			},
		},
	}
	controller.Handle(method, route, &doc, func(ctx *xx.Context) {
		p := &params{}
		err := ctx.Unmarshal(p)
		if err != nil {
			xx.HandleUnmarshalError(err, ctx)
		} else {
			roles, err := admin_model.GetRoles(p.Limit, p.Offset)
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.SendJSON(xx.MAP{"roles": roles})
			}
		}
	})
}

var CreateRole xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		Name string `required:"" desc:"角色名称"`
	}
	doc := xx.Doc {
		Title: "创建新角色",
		Params:xx.Params {
			{
				Type: xx.Query,
				Schema: &params{},
			},
		},
		Responses:xx.Responses {
			{
				Body: xx.MAP{"role": &admin_model.Role{}},
			},
		},
	}
	controller.Handle(method, route, &doc, func(ctx *xx.Context) {
		p := &params{}
		err := ctx.Unmarshal(p)
		if err != nil {
			xx.HandleUnmarshalError(err, ctx)
		} else {
			role := &admin_model.Role{Name: p.Name}
			err := role.Create()
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.SendJSON(xx.MAP{"role": role})
			}
		}
	})
}

var UpdateRole xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		ID int `required:"" desc:"角色ID"`
		Name string `required:"" desc:"新角色名称"`
	}
	doc := xx.Doc {
		Title: "更新角色名称",
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
			role := &admin_model.Role{Name: p.Name, ID: p.ID}
			err := role.Update()
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.MsgSuccess("已保存")
			}
		}
	})
}

var DeleteRoles xx.Action = func(method, route string, controller *xx.Condition) {
	type params struct {
		IDs []int `required:""`
	}
	doc := xx.Doc {
		Title: "获得角色列表",
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
			err := admin_model.DeleteRoles(p.IDs)
			if err != nil {
				ctx.MsgWarning(err.Error())
			} else {
				ctx.MsgSuccess("已删除")
			}
		}
	})
}