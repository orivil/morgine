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

var globalMiddles []*Handler

func NewGroup(tags ApiTags) *Condition {
	return DefaultServeMux.NewGroup(tags)
}

var DefaultTag = NewTagName("defaults")

var DefaultGroup = NewGroup(
	ApiTags{
		{
			Name: DefaultTag,
		},
	},
).Controller(DefaultTag)

func Use(middles ...*Handler) {
	initParser(middles...)
	globalMiddles = append(globalMiddles, middles...)
}

func initParser(hs ...*Handler) {
	for _, h := range hs {
		if h.Doc == nil {
			h.Doc = &Doc{}
		}
		if h.Doc.parser == nil {
			var err error
			h.Doc.parser, err = newParser(h.Doc.Params)
			if err != nil {
				panic(err)
			}
		}
	}
}

type Condition struct {
	middles []*Handler
	tags    ApiTags
	tagName TagName
	ApiDoc  *ApiDoc
	router  *router.Router
}

func (g *Condition) copy() *Condition {
	nc := &Condition{
		router:  g.router,
		tagName: g.tagName,
		tags:    g.tags,
		ApiDoc:  g.ApiDoc,
	}
	nc.middles = make([]*Handler, len(g.middles))
	for key, value := range g.middles {
		nc.middles[key] = value
	}
	return nc
}

func (g *Condition) Use(middles ...*Handler) *Condition {
	nc := g.copy()
	if len(middles) > 0 {
		initParser(middles...)
		nc.middles = append(nc.middles, middles...)
	}
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
	ParamType   ParamType
	Method      string
}

func (p *ParameterError) Error() string {
	return fmt.Sprintf("parameter %s is illegal: ContentType [%s], ParamType [%s], Method [%s]", p.Name, p.ContentType, p.ParamType, p.Method)
}

func (g *Condition) handle(depth int, method, route string, doc *Doc, handleFunc HandleFunc) {
	method = strings.ToUpper(method)
	if doc == nil {
		doc = &Doc{}
	}
	DocFilter(doc)
	if g.tagName == nil {
		panic("controller name is nil")
	}
	middles := append(globalMiddles, g.middles...)
	handler := &Handler {
		Doc:        doc,
		middles:    middles,
		HandleFunc: handleFunc,
	}
	initParser(handler)
	mustCheckParams(doc.parser, method)
	err := g.router.Add(method, route, handler)
	if err != nil {
		panic(err)
	}
	g.ApiDoc.add(depth+1, g.tagName, method, route, doc, middles)
}

func Handle(method, route string, doc *Doc, handleFunc HandleFunc) {
	DefaultGroup.handle(1, method, route, doc, handleFunc)
}

func mustCheckParams(pr *parser, method string) {
	for name, schema := range pr.schemas {
		typ := pr.types[name]
		ct := schema.EncodeType()
		err := &ParameterError{
			Name:        name.String(),
			ContentType: ct,
			ParamType:   typ,
			Method:      method,
		}
		switch ct {
		case param.FormDataEncodeType:
			if typ != Form {
				panic(err)
			}
		}
		switch typ {
		case Form:
			switch method {
			case http.MethodPost, http.MethodPut, http.MethodPatch:
			default:
				panic(err)
			}
		}
	}
}