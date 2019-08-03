// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

type MsgType string

const (
	MsgTypeSuccess MsgType = "success"
	MsgTypeInfo    MsgType = "info"
	MsgTypeWarning MsgType = "warning"
	MsgTypeError   MsgType = "error"
)

type Message struct {
	Type    MsgType `json:"type"`
	Content string  `json:"content"`
}
