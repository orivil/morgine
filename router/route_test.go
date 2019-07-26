// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package router_test

import (
	"github.com/orivil/morgine/router"
	"testing"
)

type dt struct {
	route   string
	prefix  string
	pattern string
}

var dts = []*dt{
	{
		route: "/",

		prefix:  "/",
		pattern: "^/$",
	},
	{
		route: "/foobar",

		prefix:  "/foobar",
		pattern: "^/foobar$",
	},
	{
		route: "/foo/bar/",

		prefix:  "/foo/bar",
		pattern: "^/foo/bar/",
	},
	{
		route: "/{foo}/{bar}",

		prefix:  "/",
		pattern: "^/(?P<foo>[^\\/^\\.]+)/(?P<bar>[^\\/^\\.]+)$",
	},
	{
		route: "/{mp}.txt",

		prefix:  "/",
		pattern: "^/(?P<mp>[^\\/^\\.]+).txt$",
	},
}

func TestInitRoute(t *testing.T) {
	for _, dt := range dts {
		prefix, pattern := router.InitRoute(dt.route)
		if dt.prefix != prefix {
			t.Errorf("prefix need: %s got: %s\n", dt.prefix, prefix)
		}
		if dt.pattern != pattern {
			t.Errorf("pattern need: %s got: %s\n", dt.pattern, pattern)
		}
	}
}

type rt struct {
	route  string
	action interface{}

	matches []*mt
}

type mt struct {
	path  string
	check func(values router.Values, action interface{}) bool
}

var rts = []*rt{
	{
		route:  "/",
		action: 1,
		matches: []*mt{
			{
				path: "/",
				check: func(values router.Values, action interface{}) bool {
					i, ok := action.(int)
					return ok && i == 1
				},
			},
			{
				path: "/foobar",
				check: func(values router.Values, action interface{}) bool {
					return action == nil
				},
			},
		},
	},
	{
		route:  "/foo/bar", // 精确匹配
		action: 2,
	},
	{
		route:  "/foo/", // 泛匹配
		action: 3,
		matches: []*mt{
			// 优先匹配长路由
			{
				path: "/foo/bar",
				check: func(values router.Values, action interface{}) bool {
					i, ok := action.(int)
					return ok && i == 2
				},
			},
			// 后匹配泛路由
			{
				path: "/foo/barbar",
				check: func(values router.Values, action interface{}) bool {
					i, ok := action.(int)
					return ok && i == 3
				},
			},
			{
				path: "/foo/bar/",
				check: func(values router.Values, action interface{}) bool {
					return action == 3
				},
			},
		},
	},
	{
		route:  "/{mp}.txt",
		action: 4,
		matches: []*mt{
			{
				path: "/123456.txt",
				check: func(values router.Values, action interface{}) bool {
					i, ok := action.(int)
					return ok && i == 4 && values().Get("mp") == "123456"
				},
			},
		},
	},
}

func TestRouter_Match(t *testing.T) {
	method := "GET"
	r := router.NewRouter()
	for _, rt := range rts {
		err := r.Add(method, rt.route, rt.action)
		if err != nil {
			t.Error(err)
		}
	}
	for _, rt := range rts {
		for _, mt := range rt.matches {
			vs, act := r.Match(method, mt.path)
			if !mt.check(vs, act) {
				t.Errorf("path [%s] need action [%v] got action [%v] values [%v]", mt.path, rt.action, act, vs())
			}
		}
	}
}
