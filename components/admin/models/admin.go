// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package models

import (
	"time"
)

type Admin struct {
	ID         int
	Username   string      `gorm:"unique_index" reg:"^[a-zA-Z0-9]{4,16}$" desc:"账号，4-16字母或数字"`
	Nickname   string      `gorm:"unique_index" desc:"昵称"`
	Password   string      `json:"-" reg:"^[a-zA-Z0-9]{8,16}$" desc:"密码，8-16字母或数字"`
	CreatedAt  *time.Time  `param:"-"`
}
