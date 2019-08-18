// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
)

type A struct {
	*B
}

type B struct {
	Name string
}

func main() {
	var a interface{} = &A{B: &B{Name: "bbb"}}
	if b, ok := a.(*B); ok {
		fmt.Println(b.Name)
	}
}
