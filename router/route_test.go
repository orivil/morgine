// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package router

import "testing"

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
		pattern: "^/foo/bar",
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
		prefix, pattern := initRoute(dt.route)
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
	path    string
	success func(action interface{}) bool
}

var rts = []*rt{
	{
		route:  "/",
		action: 1,
		matches: []*mt{
			{
				path: "/",
				success: func(action interface{}) bool {
					if action.(int) == 1 {
						return true
					}
				},
			},
		},
	},
}

func TestRouter_Match(t *testing.T) {

}
