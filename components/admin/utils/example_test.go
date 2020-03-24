// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package utils_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/orivil/morgine/components/admin/utils"
)

func ExampleWalkDirs() {
	dir := "images/avatar"
	d, err := utils.WalkDirs(dir)
	if err != nil {
		panic(err)
	}
	data, _ := formatMarshalJson(d)
	//data, _ := json.Marshal(dirs)
	fmt.Println(string(data))
	// Output:
	// {"images":{"avatar":{"admin":{},"user":{}},"cache":{},"photo":{}}}
}

func formatMarshalJson(v interface{}) (data []byte, err error) {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(v)
	if err != nil {
		return nil, err
	} else {
		return buf.Bytes(), nil
	}
}