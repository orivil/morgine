// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin

import (
	"github.com/jinzhu/gorm"
	"github.com/orivil/morgine/bundles/admin/model"
	"github.com/orivil/morgine/bundles/utils/sql"
	"github.com/orivil/morgine/cfg"
)

var DB *gorm.DB

var Register register = 0

type register int

func (r register) Init(configs cfg.Configs) {
	env := &sql.Env{}
	err := configs.Unmarshal(env)
	if err != nil {
		panic(err)
	}
	DB, err = env.Connect("admin_")
	if err != nil {
		panic(err)
	}
}

func (r register) Migrate() {
	DB.AutoMigrate(&model.Admin{})
}

func (r register) AddRoute() {

}

func (r register) RunTask() {
	panic("implement me")
}
