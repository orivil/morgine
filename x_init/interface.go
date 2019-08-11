// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package x_init

type Interface interface {
	InitConfig()
	InitDB()
	MigrateDB()
	InitRoute()
	RunTask()
}

func XInit(i ...Interface) {
	for _, it := range i {
		it.InitConfig()
	}
	for _, it := range i {
		it.InitDB()
	}
	for _, it := range i {
		it.MigrateDB()
	}
	for _, it := range i {
		it.InitRoute()
	}
	for _, it := range i {
		it.RunTask()
	}
}
