// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

type MAP map[string]interface{}

func MessageResponse(mt MsgType) *Response {
	return &Response{
		Body: msgData(mt, "some message"),
	}
}
