// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin_model

/**
1.顶级权限的path=自身的id值，次级权限的path=父级权限id值-本身权限id值
2.顶级权限的level=0,以此类推
**/

type Auth struct {
	ID int
	PID int `gorm:"index" desc:"父ID"`
	Name string `desc:"权限名"`
	C string `desc:"控制器"`
	A string `desc:"操作方法"`
	Path string `desc:"全路径"`
	Level int `desc:"权限等级"`
}
