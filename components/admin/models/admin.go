// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package models

import (
	"github.com/orivil/morgine/utils/sql"
	"time"
)

type Admin struct {
	ID int
	Username string `gorm:"unique_index"`
	Password string
	Super sql.Boolean `gorm:"index"`
	ParentID int `gorm:"index"`
	Forefather string `gorm:"index" desc:"所有祖先ID, 形如：|1|3|11|"`
	Level int `gorm:"index" desc:"账号层级，顶级管理员层级为1，随子账号逐步递增"`
	CreatedAt *time.Time
}
