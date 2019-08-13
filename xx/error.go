// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"github.com/orivil/morgine/param"
)

// 参数解析错误处理函数, 可自定义
var HandleUnmarshalError = func(err error, ctx *Context) {
	if ferr, ok := err.(*param.FieldError); ok {
		ctx.MsgWarning(ferr.Err)
	} else {
		ctx.Error(err)
	}
}
