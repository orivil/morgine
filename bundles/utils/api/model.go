// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package api

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"runtime"
)


var Models []*Model

type Model struct {
	Name   string   // 模型名称
	Desc   string   // 模型简介
	Trace  string   // 注册地址(runtime file:line)
	Fields []*Field // 模型字段
}

type Field struct {
	Name string
	Type string
	Tag  string
	Desc string
}

func NewModel(depth int, desc string, model interface{}, db *gorm.DB) *Model {
	scope := &gorm.Scope{Value: model}
	ms := scope.GetModelStruct()
	fields := ms.StructFields
	var mfs []*Field
	for _, field := range fields {
		mfs = append(mfs, &Field{
			Name: field.DBName,
			Type: field.Struct.Type.String(),
			Desc: field.Struct.Tag.Get("desc"),
			Tag:  field.Struct.Tag.Get("gorm"),
		})
	}
	return &Model {
		Name:   ms.TableName(db),
		Desc:   desc,
		Trace:  initTrace(depth),
		Fields: mfs,
	}
}

func initTrace(depth int) string {
	_, file, line, _ := runtime.Caller(depth + 1)
	return fmt.Sprintf("%s: %d", file, line)
}

func AutoMigrate(db *gorm.DB, schema interface{}, desc string) *gorm.DB {
	db.AutoMigrate(schema)
	model := NewModel(1, desc, schema, db)
	Models = append(Models, model)
	return db
}