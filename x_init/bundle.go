// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package x_init

import (
	"github.com/orivil/morgine/cfg"
)

// Bundle 接口用于规范包的初始化流程
type Bundle interface {
	// 初始化配置数据, 链接数据库等操作
	Init(configs cfg.Configs)

	// 添加路由
	AddRoute()

	// 执行定时任务, 迁移数据库等操作
	Run()
}

// Register 用于注册 bundles, 通常用于注册同一项目的 bundle
func Register(configs cfg.Configs, bs ...Bundle) {
	for _, b := range bs {
		b.Init(configs)
	}
	for _, b := range bs {
		b.AddRoute()
	}
	for _, b := range bs {
		b.Run()
	}
}

// XInit 用于初始化 bundles, 通常用于初始化不同项目的 bundle
func XInit(configs cfg.Configs, bs ...Bundle) {
	for _, b := range bs {
		b.Init(configs)
	}
}
