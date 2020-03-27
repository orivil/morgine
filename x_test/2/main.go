// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

// +build ignore

package main

import "fmt"

func main() {
	var values = map[string]map[string]string{}

	values[""] = map[string]string{}
	fmt.Println(values[""])
}
