// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"github.com/jinzhu/gorm"
	"github.com/orivil/morgine/log"
)

// TODO 测试该方法
func IsIDExist(condition *gorm.DB) bool {
	var ids []int
	condition.Order("id asc").Limit(1).Pluck("id", &ids)
	log.Error.Println("测试 IsIDExist", len(ids) > 0)
	return len(ids) > 0
}