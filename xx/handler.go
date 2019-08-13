// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"fmt"
	"github.com/orivil/morgine/param"
	"mime/multipart"
	"net/http"
	"reflect"
)

type ParamType string

const (
	Query  ParamType = "query"
	Path   ParamType = "path"
	Form   ParamType = "form"
	Header ParamType = "header"
)

type Param struct {
	Type   ParamType
	Schema interface{}
}

type Params []*Param

type Responses []*Response

type Response struct {
	Code        int
	Description string
	Headers     http.Header
	Body        interface{}
}

type Doc struct {
	Title     string
	Desc      string
	Params    Params
	Responses Responses
	parser    *parser
}

type parser struct {
	schemas   map[reflect.Type]*param.Schema
	types     map[reflect.Type]ParamType
	marshaler map[reflect.Type]marshaler
}

type marshaler func(ctx *Context) *multipart.Form

func newParser(ps Params) (par *parser, err error) {
	par = &parser{
		schemas:   make(map[reflect.Type]*param.Schema, len(ps)),
		types:     make(map[reflect.Type]ParamType, len(ps)),
		marshaler: make(map[reflect.Type]marshaler, len(ps)),
	}
	for _, p := range ps {
		schema, ok := p.Schema.(*param.Schema)
		if !ok {
			schema, err = param.NewSchema(p.Schema, nil, nil)
			if err != nil {
				return nil, err
			}
		}
		par.schemas[schema.Type] = schema
		par.types[schema.Type] = p.Type
		var m marshaler
		switch p.Type {
		case Query:
			m = func(ctx *Context) *multipart.Form {
				return &multipart.Form{Value: ctx.Query()}
			}
		case Path:
			m = func(ctx *Context) *multipart.Form {
				return &multipart.Form{Value: ctx.Path()}
			}
		case Form:
			switch schema.EncodeType() {
			case param.UrlEncodeType:
				m = func(ctx *Context) *multipart.Form {
					return &multipart.Form{Value: ctx.Form()}
				}
			case param.FormDataEncodeType:
				m = func(ctx *Context) *multipart.Form {
					return ctx.MultipartForm()
				}
			}
		case Header:
			m = func(ctx *Context) *multipart.Form {
				return &multipart.Form{Value: ctx.Request.Header}
			}
		default:
			return nil, fmt.Errorf("parameter type '%s' is not allowed", p.Type)
		}
		par.marshaler[schema.Type] = m
	}
	return par, nil
}

func (p *parser) unmarshal(vs []interface{}, ctx *Context) (err error) {
	for _, value := range vs {
		rv := reflect.ValueOf(value)
		rt := rv.Type()
		schema := p.schemas[rt]
		if schema == nil {
			return fmt.Errorf("parameter '%s' is not registered", rt)
		}
		fv := p.marshaler[rt](ctx)
		err = schema.Parse(rv.Pointer(), fv)
		if err != nil {
			return err
		}
	}
	return nil
}

type Handler struct {
	Doc        *Doc
	HandleFunc HandleFunc
	middles    []*Handler
}

type HandleFunc func(ctx *Context)

type Action func(method, route string, rg *RouteGroup)
