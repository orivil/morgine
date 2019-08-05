// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import "github.com/orivil/morgine/xx"

func main() {
	xx.Handle("GET", "/foo", func(ctx *xx.Context) {
		ctx.WriteString("bar")
	})
	xx.Handle("GET", "/{mp}.txt", func(ctx *xx.Context) {
		ctx.WriteString(ctx.Path().Get("mp"))
	})
	xx.DefaultGroup.Handle("GET", "/test", func(ctx *xx.Context) {

	})
	xx.Run()
}

func handleLogin(method, route string, c *xx.C) {

}
