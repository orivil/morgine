// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package sql

import "github.com/jinzhu/gorm"

type Boolean int

func (b Boolean) IsTrue() bool {
	return b == True
}

func GetSqlBoolean(b bool) Boolean {
	if b {
		return True
	} else {
		return False
	}
}

const (
	True Boolean = 1 + iota
	False
)

func InitOrder(field string, desc bool) (order string) {
	if field == "" {
		order = "id"
	} else {
		order = gorm.ToColumnName(field)
	}
	if desc {
		order += " desc"
	} else {
		order += " asc"
	}
	return order
}

func ToColumnName(field ...string) []string {
	for key, f := range field {
		field[key] = gorm.ToColumnName(f)
	}
	return field
}