// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"fmt"
	"github.com/orivil/morgine/router"
	"runtime"
)

func Group(tags ApiTags) *RouteGroup {
	return DefaultServeMux.Group(tags)
}

var DefaultTag = NewTagName("defaults")

var DefaultGroup = Group(
	ApiTags{
		{
			Name: DefaultTag,
			Desc: "默认路由",
		},
	},
)

type RouteGroup struct {
	middles []*Handler
	tags    ApiTags
	tagName TagName
	router  *router.Router
}

func (g *RouteGroup) copy() *RouteGroup {
	nc := &RouteGroup{
		router:  g.router,
		tagName: g.tagName,
		tags:    g.tags,
	}
	copy(nc.middles, g.middles)
	return nc
}

func (g *RouteGroup) Use(middles ...*Handler) *RouteGroup {
	nc := g.copy()
	initHandlerTrace(2, middles...)
	nc.middles = append(nc.middles, middles...)
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
	g.handle(2, method, route, handleFunc, doc)
}

func (g *RouteGroup) handle(depth int, method, route string, handleFunc HandleFunc, doc *Doc) {
	var err error
	doc.parser, err = newParser(doc.Params)
	if err != nil {
		panic(err)
	}
	if g.tagName == nil {
		g.tagName = DefaultTag
	}
	doc.tagName = g.tagName
	doc.method = method
	doc.route = route
	handler := &Handler{
		Doc:        doc,
		middles:    g.middles,
		HandleFunc: handleFunc,
	}
	initHandlerTrace(depth, handler)
	err = g.router.Add(method, route, handler)
	if err != nil {
		panic(err)
	}
}

func initHandlerTrace(depth int, h ...*Handler) {
	_, file, line, ok := runtime.Caller(depth + 1)
	if ok {
		trace := fmt.Sprintf("%s: %d", file, line)
		for _, handler := range h {
			handler.Doc.trace = trace
		}
	}
}

func Handle(method, route string, handleFunc HandleFunc) {
	DefaultGroup.handle(2, method, route, handleFunc, &Doc{})
}

var (
	Middles []*Handler
)
