// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"github.com/orivil/morgine/param"
	"github.com/orivil/morgine/router"
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
	Tag       TagName
	Method    string
	Route     string
	Params    Params
	Responses Responses
	parser    *param.Schema
}

type Middleware struct {
	Doc        *Doc
	HandleFunc HandleFunc
}

type Handler struct {
	Doc        *Doc
	Middles    []*Middleware
	HandleFunc HandleFunc
}

type HandleFunc func(ctx *Context)

type Controller struct {
	Tag      TagName
	Handlers []*Handler
}

type Condition struct {
	middles []*Middleware
	router  *router.Router
}

func (cdt *Condition) Use(middles ...*Middleware) {
	cdt.middles = append(cdt.middles, middles...)
}

func (cdt *Condition) Handle(method, route string, handleFunc HandleFunc, doc *Doc) {
	handler := &Handler{
		Doc:        doc,
		Middles:    cdt.middles,
		HandleFunc: handleFunc,
	}
	err := cdt.router.Add(method, route, handler)
	if err != nil {
		panic(err)
	}
}
