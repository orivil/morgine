// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package models

type Role struct {
	ID int
	ParentID int `gorm:"index"`
	Name string
}
