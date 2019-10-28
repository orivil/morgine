// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import "github.com/orivil/morgine/param"

type MAP map[string]interface{}

// 错误处理函数, 可自定义
var HandleError = func(ctx *Context, err error) {
	if verr, ok := err.(*param.ValidatorErr); ok {
		HandleMessage(ctx, MsgTypeError, verr.Field + ":" + verr.Message)
	} else {
		_ = ctx.TraceError(2, err)
	}
}

// 消息处理函数, 可自定义
var HandleMessage = func(ctx *Context, mt MsgType, msg string) {
	_ = ctx.SendJSON(Message{Type:mt, Content: msg})
}