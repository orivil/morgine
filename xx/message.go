// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

type MsgType string

const (
	MsgSuccess MsgType = "success"
	MsgInfo    MsgType = "info"
	MsgWarning MsgType = "warning"
	MsgError   MsgType = "error"
)

type Message struct {
	Type    MsgType `json:"type" xml:"type"`
	Content string  `json:"content" xml:"content"`
}