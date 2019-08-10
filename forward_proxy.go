// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package main

import (
	"fmt"
)

var passwd = "wen123456"

func main() {
	var s = []string{"null"}
	fmt.Println(len(s))
	if len(s) == 1 && s[0] == "" || s[0] == "null" {
		fmt.Println(true)
	}
	fmt.Println(s)
}
