// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import "github.com/orivil/morgine/param"

type MAP map[string]interface{}

func MessageResponse(mt MsgType) *Response {
	return &Response{
		Body: msgData(mt, "some message"),
	}
}

func MessageData(mt MsgType, msg string) map[string]*Message {
	return msgData(mt, msg)
}

// 参数解析错误处理函数, 可自定义
var HandleUnmarshalError = func(err error, ctx *Context) {
	if verr, ok := err.(*param.ValidatorErr); ok {
		ctx.MsgWarning(verr.Message)
	} else {
		ctx.Error(err)
	}
}
