// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package db

import (
	"github.com/jinzhu/gorm"
	"github.com/orivil/morgine/bundles/utils/sql"
)

var GORM *gorm.DB

func InitDB() {
	err := sql.InitConfig("admin.yml", func(db *gorm.DB) {
		GORM = db
	})
	if err != nil {
		panic(err)
	}
}
