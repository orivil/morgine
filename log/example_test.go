// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package log_test

import (
	"github.com/orivil/morgine/log"
	"os"
)

func ExampleSetOutput() {

	log.Info.SetFlags(0)

	// 自定义写入接口
	log.SetOutput(log.FlagInfo|log.FlagDanger, os.Stdout)

	// 打印一行 info 信息
	log.Info.Println("hello")

	// Output:
	// [info] hello
}
