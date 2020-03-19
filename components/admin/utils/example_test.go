// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func ExampleWalkDirs() {
	dir := "../../../"
	dirs, err := WalkDirs(dir)
	if err != nil {
		panic(err)
	}
	buf := bytes.NewBuffer(nil)
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(dirs)
	if err != nil {
		panic(err)
	}

	fmt.Println(buf.String())
	// Output:
	// 134
}
