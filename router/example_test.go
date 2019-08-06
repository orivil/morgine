// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package router_test

import (
	"fmt"
	"github.com/orivil/morgine/router"
)

func ExampleNewRouter() {
	r := router.NewRouter()
	err := r.Add("GET", "/{mp}.txt", 1)
	if err != nil {
		// if err is not nil, it must be regular expresion error
		panic(err)
	}

	// match all route
	r.Add("GET", "/", 2)

	values, action := r.Match("GET", "/123456.txt")
	fmt.Println(action.(int) == 1)
	fmt.Println(values().Get("mp") == "123456")
	_, action = r.Match("GET", "/foo/bar")
	fmt.Println(action.(int) == 2)
	_, action = r.Match("GET", "/")
	fmt.Println(action.(int) == 2)

	// Output:
	// true
	// true
	// true
	// true
}
