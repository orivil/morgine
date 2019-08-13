// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package x_init

import "github.com/orivil/morgine/cfg"

type Interface interface {
	Init(configs cfg.Configs)
	Migrate()
	AddRoute()
	RunTask()
}

func XInit(configs cfg.Configs, i ...Interface) {
	for _, it := range i {
		it.Init(configs)
	}
	for _, it := range i {
		it.Migrate()
	}
	for _, it := range i {
		it.AddRoute()
	}
	for _, it := range i {
		it.RunTask()
	}
}
