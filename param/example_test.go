// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package param_test

import (
	"fmt"
	"github.com/orivil/morgine/param"
	"mime/multipart"
	"net/textproto"
	"time"
	"unsafe"
)

type embed struct {

	// 枚举限制, 包括 int, string 等类型
	Enum string `enum:"1 2 3"`

	// email 格式限制
	Email string `email:""`
}

type Param struct {
	// required 限制, 字符串长度限制, 正则匹配
	Str string `required:"" len:"2-4" reg:"[\\w]+"`

	// 数量限制, 包括 []int, []float64 等类型
	Strs []string `item:"2-4"`

	embed // 嵌入结构体

	// 数字大小限制, 包括 float32 float64 类型
	Num int `num:"1<x<=3"`

	// 参数重命名
	Bool bool `param:"bool"`

	// 忽略字段
	IgnoreMe string `param:"-"`

	// 文件大小限制(单位: KB)
	Photo param.FileHandler `size:"20-300"`

	// 文件后缀名限制
	Image1 param.FileHandler `exts:".jpg .jpeg .png .gif"`

	// MIME type 具体类型限制
	Image2 param.FileHandler `mime:"image/jpeg image/png"`

	// MIME type 主类型限制
	Image3 param.FileHandler `mime:"image"`

	// 文件数量限制
	Files param.FileHandler `item:"1-3"`

	// 时间格式, 默认解析 RFC3339: "2006-01-02T15:04:05" 格式时间, 如果要解析其他格式, 则定义一个 time-layout 标签
	Birthday *time.Time `time-layout:"2006-01-02T15:04:05"`
}

// 模拟数据源
var formData = &multipart.Form{
	// 模拟字段数据
	Value: map[string][]string{
		"Str":      {"12"},
		"Strs":     {"1", "2"},
		"Enum":     {"2"},
		"Email":    {"author@example.com"},
		"Num":      {"2"},
		"bool":     {},
		"Birthday": {"2006-01-02T5:04:05"},
	},
	// 模拟文件数据
	File: map[string][]*multipart.FileHeader{
		// 模拟 20KB 大小文件
		"Photo": {&multipart.FileHeader{Filename: "a.jpg", Size: 20 << 10}},
		// 模拟 "image/jpeg" 文件
		"Image3": {&multipart.FileHeader{Filename: "a.jpg", Header: textproto.MIMEHeader{"Content-Type": []string{"image/jpeg"}}}},
		// 模拟 4 个上传文件
		"Files": {
			&multipart.FileHeader{Filename: "a.jpg"},
			&multipart.FileHeader{Filename: "b.jpg"},
			&multipart.FileHeader{Filename: "c.jpg"},
			&multipart.FileHeader{Filename: "d.jpg"},
		},
	},
}

func ExampleSchema() {

	// 新建解析器
	var parser, err = param.NewSchema(&Param{}, nil, nil)
	if err != nil {
		panic(err)
	}

	// 需要解析的参数
	p := &Param{
		// 定义文件处理函数
		Photo:  fileHandler,
		Image1: fileHandler,
		Image2: fileHandler,
		Image3: fileHandler,
		Files:  fileHandler,
	}

	// 验证并解析参数
	err = parser.Parse(uintptr(unsafe.Pointer(p)), formData)
	if err != nil {
		//fmt.Println(err)

		if ce, ok := err.(*param.ValidatorErr); ok { // 可获得具体的条件数据
			fmt.Printf("%s: %s\n", ce.Field, ce.Message)
		}
	}
	// Output:
	// Files: 1-3 items
}

// 文件存储函数
var fileHandler param.FileHandler = func(field string, header *multipart.FileHeader) error {
	return nil
}
