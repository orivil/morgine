// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"github.com/orivil/morgine/router"
)

func NewGroup(tags ApiTags) *RouteGroup {
	return DefaultServeMux.NewGroup(tags)
}

var DefaultTag = NewTagName("defaults")

var DefaultGroup = NewGroup(
	ApiTags{
		{
			Name: DefaultTag,
			Desc: "默认路由",
		},
	},
)

func Use(middles ...*Handler) {
	DefaultGroup = DefaultGroup.Use(middles...)
}

type RouteGroup struct {
	middles []*Handler
	tags    ApiTags
	tagName TagName
	apiDoc  *apiDoc
	router  *router.Router
}

func (g *RouteGroup) copy() *RouteGroup {
	nc := &RouteGroup{
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

func (g *RouteGroup) Use(middles ...*Handler) *RouteGroup {
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

func (g *RouteGroup) Controller(name TagName) *RouteGroup {
	if !g.tags.checkIsSubTag(name) {
		panic("need the sub of the initialized tags")
	}
	nc := g.copy()
	nc.tagName = name
	return nc
}

func (g *RouteGroup) Handle(method, route string, handleFunc HandleFunc, doc *Doc) {
	g.handle(1, method, route, handleFunc, doc)
}

func (g *RouteGroup) handle(depth int, method, route string, handleFunc HandleFunc, doc *Doc) {
	if doc == nil {
		doc = &Doc{}
	}
	var err error
	doc.parser, err = newParser(doc.Params)
	if err != nil {
		panic(err)
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

func Handle(method, route string, handleFunc HandleFunc, doc *Doc) {
	DefaultGroup.handle(1, method, route, handleFunc, doc)
}
