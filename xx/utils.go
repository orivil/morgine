// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"github.com/orivil/morgine/param"
)

type MAP map[string]interface{}

// 错误处理函数, 可自定义
var HandleError = func(ctx *Context, err error) {
	if verr, ok := err.(*param.ValidatorErr); ok {
		SendMessage(ctx, MsgTypeError, verr.Field + ":" + verr.Message)
	} else {
		_ = ctx.TraceError(2, err)
	}
}

// 消息处理函数, 可自定义
var SendMessage = func(ctx *Context, mt MsgType, msg string) {
	_ = ctx.SendJSON(Message{Type:mt, Content: msg})
}

// json 数据处理函数，可自定义
var SendJson = func(ctx *Context, code StatusCode, data interface{}) {
	_ = ctx.SendJSON(JsonData(code, data))
}

// 消息数据模型
var MessageData = func(mt MsgType, msg string) *Message {
	return &Message {
		Type: mt,
		Content: msg,
	}
}

// json 数据结构函数，可自定义
var JsonData = func(code StatusCode, data interface{}) MAP {
	return MAP {
		"code": code,
		"data": data,
	}
}

// http.Error() 返回结果模型
var HttpErrorResponse = func(err string, statusCode int) *Response {
	return &Response {
		Code:        statusCode,
		Description: "",
		Body:        err,
	}
}