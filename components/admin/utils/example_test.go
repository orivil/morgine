// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package utils_test

import (
	"encoding/json"
	"fmt"
	"github.com/orivil/morgine/components/admin/utils"
)

func ExampleWalkDirs() {
	dir := "images"
	dirs, err := utils.WalkDirs(dir)
	if err != nil {
		panic(err)
	}
	data, _ := json.Marshal(dirs)
	fmt.Println(string(data))
	// Output:
	// {"images":{"avatar":{"admin":{},"user":{}},"cache":{},"photo":{}}}
}
