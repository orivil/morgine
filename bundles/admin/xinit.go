// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin

import (
	"github.com/orivil/morgine/bundles/admin/db"
	"github.com/orivil/morgine/bundles/admin/model"
)

type Register int

func (r Register) InitConfig() {

}

func (r Register) InitDB() {
	db.InitDB()
}

func (r Register) MigrateDB() {
	db.GORM.AutoMigrate(&model.Admin{})
}

func (r Register) InitRoute() {

}

func (r Register) RunTask() {

}
