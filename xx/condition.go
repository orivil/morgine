// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"fmt"
	"github.com/orivil/morgine/param"
	"github.com/orivil/morgine/router"
	"net/http"
	"strings"
)

func NewGroup(tags ApiTags) *Condition {
	return DefaultServeMux.NewGroup(tags)
}

var DefaultTag = NewTagName("defaults")

var DefaultCondition = NewGroup(
	ApiTags{
		{
			Name: DefaultTag,
			Desc: "默认路由",
		},
	},
)

func Use(middles ...*Handler) {
	DefaultCondition = DefaultCondition.Use(middles...)
}

type Condition struct {
	middles []*Handler
	tags    ApiTags
	tagName TagName
	apiDoc  *apiDoc
	router  *router.Router
}

func (g *Condition) copy() *Condition {
	nc := &Condition{
		router:  g.router,
		tagName: g.tagName,
		tags:    g.tags,
		apiDoc:  g.apiDoc,
	}
	nc.middles = make([]*Handler, len(g.middles))
	for key, value := range g.middles {
		nc.middles[key] = value
	}
	return nc
}

func (g *Condition) Use(middles ...*Handler) *Condition {
	nc := g.copy()
	for _, middle := range middles {
		if middle.Doc == nil {
			middle.Doc = &Doc{}
		}
		if middle.Doc.parser == nil {
			var err error
			middle.Doc.parser, err = newParser(middle.Doc.Params)
			if err != nil {
				panic(err)
			}
		}
	}
	nc.middles = append(nc.middles, middles...)
	nc.apiDoc.use(middles...)
	return nc
}

func (g *Condition) Controller(name TagName) *Condition {
	if !g.tags.checkIsSubTag(name) {
		panic("need the sub of the initialized tags")
	}
	if !g.tags.checkIsEndTag(name) {
		panic("need the leaf node of the initialized tags")
	}
	nc := g.copy()
	nc.tagName = name
	return nc
}

func (g *Condition) Handle(method, route string, doc *Doc, handleFunc HandleFunc) {
	g.handle(1, method, route, doc, handleFunc)
}

type ParameterError struct {
	Name        string
	ContentType param.EncodeType
	Type        ParamType
	Method      string
}

func (p *ParameterError) Error() string {
	return fmt.Sprintf("parameter %s is illegal: ContentType [%s], Location [%s], Method [%s]", p.Name, p.ContentType, p.Type, p.Method)
}

func (g *Condition) handle(depth int, method, route string, doc *Doc, handleFunc HandleFunc) {
	method = strings.ToUpper(method)
	if doc == nil {
		doc = &Doc{}
	}
	var err error
	doc.parser, err = newParser(doc.Params)
	if err != nil {
		panic(err)
	}
	for name, schema := range doc.parser.schemas {
		switch ct := schema.EncodeType(); ct {
		case param.FormDataEncodeType:
			typ := doc.parser.types[name]
			switch typ {
			case Form:
			default:
				panic(&ParameterError{
					Name:        name.String(),
					ContentType: ct,
					Type:        typ,
					Method:      method,
				})
			}
			switch method {
			case http.MethodPost, http.MethodPut:
			default:
				panic(&ParameterError{
					Name:        name.String(),
					ContentType: ct,
					Type:        typ,
					Method:      method,
				})
			}
		}
	}
	if g.tagName == nil {
		g.tagName = DefaultTag
	}
	handler := &Handler{
		Doc:        doc,
		middles:    g.middles,
		HandleFunc: handleFunc,
	}
	err = g.router.Add(method, route, handler)
	if err != nil {
		panic(err)
	}
	g.apiDoc.handle(depth+1, g.tagName, method, route, doc, g.middles)
}

func Handle(method, route string, doc *Doc, handleFunc HandleFunc) {
	DefaultCondition.handle(1, method, route, doc, handleFunc)
}
