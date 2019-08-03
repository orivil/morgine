// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"fmt"
	"github.com/orivil/morgine/log"
	"github.com/orivil/morgine/router"
	"github.com/orivil/morgine/utils/ip"
	"net/http"
	"runtime"
	"sync"
)

type ServeMux struct {
	r               *router.Router
	NotFoundHandler http.HandlerFunc
}

var contextPool = sync.Pool{
	New: func() interface{} {
		return &Context{}
	},
}

func (mux *ServeMux) Group(tags ApiTags) *group {
	return &group{
		tags:   tags,
		router: mux.r,
	}
}

func (mux *ServeMux) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	vs, act := mux.r.Match(req.Method, req.URL.Path)
	if act != nil {
		ctx := contextPool.Get().(*Context)
		defer func() {
			contextPool.Put(ctx)
			err := recover()
			if err != nil && err != http.ErrAbortHandler {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Panic.Printf("%s\n http panic: %s\n %s", GetRequestInfo(req), err, buf)
				http.Error(ctx.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		ctx = initContext(ctx, writer, req, vs, act.(*Handler))
		ctx.handle()
	} else {
		mux.NotFoundHandler(writer, req)
	}
}

func GetRequestInfo(r *http.Request) string {
	return fmt.Sprintf("| %16s | %8s | %s", ip.GetHttpRequestIP(r), r.Method, r.Host+r.URL.Path)
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		r:               router.NewRouter(),
		NotFoundHandler: http.NotFound,
	}
}

var DefaultServeMux = NewServeMux()

func Group(tags ApiTags) *group {
	return DefaultServeMux.Group(tags)
}

type group struct {
	middles []*Handler
	tags    ApiTags
	tagName TagName
	router  *router.Router
}

func (g *group) copy() *group {
	nc := &group{
		router:  g.router,
		tagName: g.tagName,
		tags:    g.tags,
	}
	copy(nc.middles, g.middles)
	return nc
}

func (g *group) Use(middles ...*Handler) *group {
	nc := g.copy()
	nc.middles = append(nc.middles, middles...)
	return nc
}

func (g *group) Controller(name TagName) *group {
	if !g.tags.checkIsSubTag(name) {
		panic("need the sub of the initialized tags")
	}
	nc := g.copy()
	nc.tagName = name
	return nc
}

func (g *group) Handle(method, route string, handleFunc HandleFunc, doc *Doc) {
	var err error
	doc.parser, err = newParser(doc.Params)
	if err != nil {
		panic(err)
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
}
