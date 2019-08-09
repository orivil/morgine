// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package rand_test

import (
	"fmt"
	"github.com/orivil/morgine/utils/rand"
)

func ExampleNewUUID() {

	// 生成指定位数的随机字符串
	str := rand.NewRandByte(64)
	fmt.Println(string(str))
}
