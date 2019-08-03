// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.
package log

import (
	"io"
	"log"
	"os"
)

type flag int

const (
	FlagInit flag = 1 << iota
	FlagInfo
	FlagError
	FlagWarning
	FlagDanger
	FlagEmergency
	FlagPanic
)

// Init 用于程序初始化时保存一些重要信息
var Init = log.New(os.Stdout, " ", log.LstdFlags)

var Info = log.New(os.Stdout, "[info] ", log.LstdFlags)

var Warning = log.New(os.Stderr, "[warning] ", log.LstdFlags|log.Llongfile)

var Danger = log.New(os.Stderr, "[danger] ", log.LstdFlags|log.Llongfile)

var Error = log.New(os.Stderr, "[error] ", log.LstdFlags|log.Llongfile)

var Emergency = log.New(os.Stderr, "[emergency] ", log.LstdFlags|log.Llongfile)

// 用于记录 panic 信息, 并非触发 panic
var Panic = log.New(os.Stderr, "[panic] ", log.LstdFlags|log.Llongfile)

func SetOutput(flag flag, writer io.Writer) {
	if flag&FlagInit != 0 {
		Init.SetOutput(writer)
	}
	if flag&FlagInfo != 0 {
		Info.SetOutput(writer)
	}
	if flag&FlagWarning != 0 {
		Warning.SetOutput(writer)
	}
	if flag&FlagError != 0 {
		Error.SetOutput(writer)
	}
	if flag&FlagDanger != 0 {
		Danger.SetOutput(writer)
	}
	if flag&FlagEmergency != 0 {
		Emergency.SetOutput(writer)
	}
	if flag&FlagPanic != 0 {
		Panic.SetOutput(writer)
	}
}
