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

type Responses []Response

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
	schemas map[string]*param.Schema
	types   map[string]ParamType
}

func newParser(ps Params) (par *parser, err error) {
	par = &parser{
		schemas: make(map[string]*param.Schema, len(ps)),
		types:   make(map[string]ParamType, len(ps)),
	}
	for _, p := range ps {
		schema, ok := p.Schema.(*param.Schema)
		if !ok {
			schema, err = param.NewSchema(p.Schema, nil, nil)
			if err != nil {
				return nil, err
			}
		}
		par.schemas[schema.Name] = schema
		par.types[schema.Name] = p.Type
	}
	return par, nil
}

func (p *parser) unmarshal(vs []interface{}, ctx *Context) (err error) {
	for _, value := range vs {
		rv := reflect.ValueOf(value)
		name := rv.Type().Name()
		schema := p.schemas[name]
		if schema == nil {
			return fmt.Errorf("parameter '%s' is not registered", name)
		}
		t := p.types[name]
		var fv *multipart.Form
		switch t {
		case Query:
			fv = &multipart.Form{Value: ctx.Query()}
		case Path:
			fv = &multipart.Form{Value: ctx.Path()}
		case Form:
			switch schema.EncodeType() {
			case param.UrlEncodeType:
				fv = &multipart.Form{Value: ctx.Form()}
			case param.FormDataEncodeType:
				fv, err = ctx.parseMultipartForm()
				if err != nil {
					return err
				}
			}
		case Header:
			fv = &multipart.Form{Value: ctx.Request.Header}
		default:
			return fmt.Errorf("parameter type '%s' is not allowed", t)
		}
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
