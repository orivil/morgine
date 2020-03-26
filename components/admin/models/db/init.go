// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package db

import (
	"github.com/jinzhu/gorm"
	"github.com/orivil/morgine/components/admin/models"
)

var (
	DB         *gorm.DB
)

func Init(gdb *gorm.DB) {
	DB = gdb
	DB.AutoMigrate (
		&models.Admin{},
		&models.SuperAdmin{},
	)
}
