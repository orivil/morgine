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
	"time"
)

type RequestLogger func(req *http.Request, costTime time.Duration, statusCode int)

type ServeMux struct {
	r               *router.Router
	ErrHandler      func(w http.ResponseWriter, error string, code int)
	RequestLogger   RequestLogger
	apiDoc          *apiDoc
}

func NewServeMux(r *router.Router) *ServeMux {
	return &ServeMux {
		r:               r,
		ErrHandler:      http.Error,
		RequestLogger: func(req *http.Request, costTime time.Duration, statusCode int) {
			log.Info.Printf("| %14s | %4d %s \n\n", costTime, statusCode, GetRequestInfo(req))
		},
		apiDoc: newApiDoc(),
	}
}

var contextPool = sync.Pool {
	New: func() interface{} {
		return &Context{}
	},
}

func (mux *ServeMux) ApiDoc() *apiDoc {
	return mux.apiDoc
}

func (mux *ServeMux) NewGroup(tags ApiTags) *Condition {
	mux.apiDoc.Tags = append(mux.apiDoc.Tags, tags...)
	return &Condition{
		tags:   tags,
		router: mux.r,
		apiDoc: mux.apiDoc,
	}
}

func (mux *ServeMux) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	if mux.RequestLogger != nil {
		res := &response{ResponseWriter: writer}
		writer = res
		start := time.Now()
		defer func() {
			cost := time.Since(start)
			mux.RequestLogger(req, cost, res.statusCode)
		}()
	}
	vs, act := mux.r.Match(req.Method, req.URL.Path)
	if act != nil {
		ctx := contextPool.Get().(*Context)
		defer func() {
			err := recover()
			if err != nil && err != http.ErrAbortHandler {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Panic.Printf("%s\n http panic: %s\n %s", GetRequestInfo(req), err, buf)
				mux.ErrHandler(ctx.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			} else if ctx.err != nil {
				mux.ErrHandler(ctx.Writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
			contextPool.Put(ctx)
		}()
		ctx = initContext(ctx, writer, req, vs, act.(*Handler))
		ctx.handle()
	} else {
		NotFoundHandler.ServeHTTP(writer, req)
	}
}

func GetRequestInfo(r *http.Request) string {
	return fmt.Sprintf("| %16s | %8s | %s", ip.GetHttpRequestIP(r), r.Method, r.Host+r.URL.Path)
}

var DefaultServeMux = NewServeMux(router.NewRouter())