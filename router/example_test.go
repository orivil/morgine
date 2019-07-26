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
		// must be regular expresion error
		panic(err)
	}
	values, action := r.Match("GET", "/123456.txt")

	fmt.Println(action.(int) == 1)
	fmt.Println(values().Get("mp") == "123456")

	// Output:
	// true
	// true
}
