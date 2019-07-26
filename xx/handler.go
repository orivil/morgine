// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"github.com/orivil/morgine/param"
	"net/http"
)

type ParamType string

const (
	Query  ParamType = "query"
	Form   ParamType = "form"
	Body   ParamType = "body"
	Header ParamType = "header"
)

type Param struct {
	Type   ParamType
	Schema interface{}
}

type Params []*Param

type Responses []Response

type Response struct {
	Code        int
	Description string
	Headers     http.Header
	Body        interface{}
}

type Doc struct {
	Method    string
	Route     string
	Params    Params
	Responses Responses
	parser    param.Parser
}

type Handler struct {
	Doc    Doc
	Handle func(ctx *Context)
}

type Tag *string

func TagName(name string) Tag {
	return &name
}

type Tags []Tag

type Controller struct {
	Tag     Tag
	Actions []*Handler
}
