// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package param

import (
	"fmt"
	"mime/multipart"
	"net/textproto"
	"reflect"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

type param struct {
	Str          string `param:"str"`
	Int          int
	Int32        int32
	Int64        int64
	Float32      float32
	Float64      float64
	Bool         bool
	SliceStr     []string
	SliceInt     []int
	SliceInt32   []int32
	SliceInt64   []int64
	SliceFloat32 []float32
	SliceFloat64 []float64
	SliceBool    []bool
	File         FileHandler
	Time         *time.Time
	tag          *tag `param:"-"`
}

func initValidator(p *param) *Validator {
	if p.tag == nil || p.tag.field == "" {
		return nil
	}
	// 独立检验单个字段，避免数据被干扰
	v := NewValidator(p)
	var c *condition
	switch p.tag.field {
	case "Str":
		c = v.Field(&p.Str)
	case "Int":
		c = v.Field(&p.Int)
	case "Int32":
		c = v.Field(&p.Int32)
	case "Int64":
		c = v.Field(&p.Int64)
	case "Float32":
		c = v.Field(&p.Float32)
	case "Float64":
		c = v.Field(&p.Float64)
	case "Bool":
		c = v.Field(&p.Bool)
	case "SliceStr":
		c = v.Field(&p.SliceStr)
	case "SliceInt":
		c = v.Field(&p.SliceInt)
	case "SliceInt32":
		c = v.Field(&p.SliceInt32)
	case "SliceInt64":
		c = v.Field(&p.SliceInt64)
	case "SliceFloat32":
		c = v.Field(&p.SliceFloat32)
	case "SliceFloat64":
		c = v.Field(&p.SliceFloat64)
	case "SliceBool":
		c = v.Field(&p.SliceBool)
	case "File":
		c = v.Field(&p.File)
	case "Time":
		c = v.Field(&p.Time)
	default:
		panic(fmt.Sprintf("field %s not exist", p.tag.field))
	}
	err := c.Syntax(p.tag.tag)
	if err != nil {
		panic(err)
	}
	return v
}

type tag struct {
	field  string
	tag    string
	msg    string
	key    string // value key, default is the "field"
	values []string
	files  []multipart.FileHeader
	line   string
}

type value struct {
	field string
	form  []string
	check func(p *param) bool
}

var values = []value{
	{
		field: "str",
		form:  []string{"str"},
		check: func(p *param) bool {
			return p.Str == "str"
		},
	},

	{
		field: "Int",
		form:  []string{"2"},
		check: func(p *param) bool {
			return p.Int == 2
		},
	},

	{
		field: "Int32",
		form:  []string{"2"},
		check: func(p *param) bool {
			return p.Int32 == 2
		},
	},

	{
		field: "Int64",
		form:  []string{"2"},
		check: func(p *param) bool {
			return p.Int64 == 2
		},
	},

	{
		field: "Float32",
		form:  []string{"2"},
		check: func(p *param) bool {
			return p.Float32 == 2
		},
	},

	{
		field: "Float64",
		form:  []string{"2"},
		check: func(p *param) bool {
			return p.Float64 == 2
		},
	},

	{
		field: "Bool",
		form:  []string{"1"},
		check: func(p *param) bool {
			return p.Bool == true
		},
	},

	{
		field: "SliceStr",
		form:  []string{"1", "2"},
		check: func(p *param) bool {
			return reflect.DeepEqual(p.SliceStr, []string{"1", "2"})
		},
	},

	{
		field: "SliceInt",
		form:  []string{"1", "2"},
		check: func(p *param) bool {
			return reflect.DeepEqual(p.SliceInt, []int{1, 2})
		},
	},
	{
		field: "SliceInt32",
		form:  []string{"1", "2"},
		check: func(p *param) bool {
			return reflect.DeepEqual(p.SliceInt32, []int32{1, 2})
		},
	},
	{
		field: "SliceInt64",
		form:  []string{"1", "2"},
		check: func(p *param) bool {
			return reflect.DeepEqual(p.SliceInt64, []int64{1, 2})
		},
	},
	{
		field: "SliceFloat32",
		form:  []string{"1", "2"},
		check: func(p *param) bool {
			return reflect.DeepEqual(p.SliceFloat32, []float32{1, 2})
		},
	},
	{
		field: "SliceFloat64",
		form:  []string{"1", "2"},
		check: func(p *param) bool {
			return reflect.DeepEqual(p.SliceFloat64, []float64{1, 2})
		},
	},
	{
		field: "SliceBool",
		form:  []string{"1", "1"},
		check: func(p *param) bool {
			return reflect.DeepEqual(p.SliceBool, []bool{true, true})
		},
	},
	{
		field: "Time",
		form:  []string{"2018-12-20T00:00:00"},
		check: func(p *param) bool {
			t, err := time.Parse(DefaultTimeLayout, "2018-12-20T00:00:00")
			if err != nil {
				panic(err)
			}
			return p.Time.Equal(t)
		},
	},
}

//var now = time.Now()

// 测试解析值
func TestParser_Parse(t *testing.T) {
	p := &param{}
	parser, err := NewSchema(p, initValidator(p), nil)
	if err != nil {
		t.Error(err)
	} else {
		for _, value := range values {
			p := &param{
				Time: &time.Time{},
			}
			var form multipart.Form
			form.Value = map[string][]string{value.field: value.form}
			err := parser.Parse(uintptr(unsafe.Pointer(p)), &form)
			if err != nil {
				t.Fatalf("field %s parse got error: %v", value.field, err)
			}
			if !value.check(p) {
				t.Fatalf("field %s check failed\n", value.field)
			}
		}
	}
}

var tags = []tag{
	// test required
	{line: getFileLine(), field: "Str", tag: `required:"required"`, msg: "required", key: "str", values: []string{""}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `required:"required"`, msg: "", key: "str", values: []string{"value"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int", tag: `required:"required"`, msg: "required", values: []string{""}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int32", tag: `required:"required"`, msg: "required", values: []string{""}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int32", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int64", tag: `required:"required"`, msg: "required", values: []string{""}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int64", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Float32", tag: `required:"required"`, msg: "required", values: []string{""}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Float32", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Float64", tag: `required:"required"`, msg: "required", values: []string{""}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Float64", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceStr", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceStr", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceInt", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceInt32", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt32", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceInt64", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt64", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceFloat32", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceFloat32", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceFloat64", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceFloat64", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceBool", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceBool", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "File", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "File", tag: `required:"required"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{{Filename: "pohot.jpg"}}},

	{line: getFileLine(), field: "Time", tag: `required:"required"`, msg: "required", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Time", tag: `required:"required"`, msg: "", values: []string{"2006-01-02T15:04:05"}, files: []multipart.FileHeader{}},

	// test len
	{line: getFileLine(), field: "Str", tag: `len:"2-4" len-msg:"2-4"`, msg: "2-4", key: "str", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `len:"2-4" len-msg:"2-4"`, msg: "", key: "str", values: []string{"12"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `len:"2-4" len-msg:"2-4"`, msg: "2-4", key: "str", values: []string{"12345"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `len:"2-4" len-msg:"2-4"`, msg: "2-4", key: "str", values: []string{"一"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `len:"2-4" len-msg:"2-4"`, msg: "", key: "str", values: []string{"一二"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `len:"2-4" len-msg:"2-4"`, msg: "2-4", key: "str", values: []string{"一二三四五"}, files: []multipart.FileHeader{}},

	// test regular expression
	{line: getFileLine(), field: "Str", tag: `reg:"^[\\w]+$" reg-msg:"msg"`, msg: "", key: "str", values: []string{"12345"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `reg:"^[\\w]+$" reg-msg:"msg"`, msg: "msg", key: "str", values: []string{"123;45"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `email:"msg"`, msg: "", key: "str", values: []string{"admin@orivil.com"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `email:"msg"`, msg: "msg", key: "str", values: []string{"admin@orivilcom"}, files: []multipart.FileHeader{}},

	// test enum
	{line: getFileLine(), field: "Str", tag: `enum:"1 2 3 4" enum-msg:"msg"`, msg: "", key: "str", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Str", tag: `enum:"1 2 3 4" enum-msg:"msg"`, msg: "msg", key: "str", values: []string{"0"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int", tag: `enum:"1 2 3 4" enum-msg:"msg"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int", tag: `enum:"1 2 3 4" enum-msg:"msg"`, msg: "msg", values: []string{"0"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int32", tag: `enum:"1 2 3 4" enum-msg:"msg"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int32", tag: `enum:"1 2 3 4" enum-msg:"msg"`, msg: "msg", values: []string{"0"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int64", tag: `enum:"1 2 3 4" enum-msg:"msg"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int64", tag: `enum:"1 2 3 4" enum-msg:"msg"`, msg: "msg", values: []string{"0"}, files: []multipart.FileHeader{}},

	// test num
	{line: getFileLine(), field: "Int", tag: `num:"2<=x<4" num-msg:"msg"`, msg: "msg", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int", tag: `num:"2<=x<4" num-msg:"msg"`, msg: "", values: []string{"2"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int", tag: `num:"2<=x<4" num-msg:"msg"`, msg: "", values: []string{"3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int", tag: `num:"2<=x<4" num-msg:"msg"`, msg: "msg", values: []string{"4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "msg", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "msg", values: []string{"2"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "", values: []string{"3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "", values: []string{"4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int32", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "msg", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int32", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "msg", values: []string{"2"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int32", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "", values: []string{"3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int32", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "", values: []string{"4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Int64", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "msg", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int64", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "msg", values: []string{"2"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int64", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "", values: []string{"3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Int64", tag: `num:"2<x<=4" num-msg:"msg"`, msg: "", values: []string{"4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Float32", tag: `num:"2.0<x<=4.0" num-msg:"msg"`, msg: "msg", values: []string{"1.0"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Float32", tag: `num:"2.0<x<=4.0" num-msg:"msg"`, msg: "msg", values: []string{"2.0"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Float32", tag: `num:"2.0<x<=4.0" num-msg:"msg"`, msg: "", values: []string{"3.0"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Float32", tag: `num:"2.0<x<=4.0" num-msg:"msg"`, msg: "", values: []string{"4.0"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "Float64", tag: `num:"2.0<x<=4.0" num-msg:"msg"`, msg: "msg", values: []string{"1.0"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Float64", tag: `num:"2.0<x<=4.0" num-msg:"msg"`, msg: "msg", values: []string{"2.0"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Float64", tag: `num:"2.0<x<=4.0" num-msg:"msg"`, msg: "", values: []string{"3.0"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "Float64", tag: `num:"2.0<x<=4.0" num-msg:"msg"`, msg: "", values: []string{"4.0"}, files: []multipart.FileHeader{}},

	// test item
	{line: getFileLine(), field: "SliceStr", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceStr", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1", "2", "3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceStr", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{}}, // not required
	{line: getFileLine(), field: "SliceStr", tag: `item:"1-3" item-msg:"1-3"`, msg: "1-3", values: []string{"1", "2", "3", "4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceInt", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1", "2", "3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{}}, // not required
	{line: getFileLine(), field: "SliceInt", tag: `item:"1-3" item-msg:"1-3"`, msg: "1-3", values: []string{"1", "2", "3", "4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceInt32", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt32", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1", "2", "3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt32", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{}}, // not required
	{line: getFileLine(), field: "SliceInt32", tag: `item:"1-3" item-msg:"1-3"`, msg: "1-3", values: []string{"1", "2", "3", "4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceInt64", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt64", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1", "2", "3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceInt64", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{}}, // not required
	{line: getFileLine(), field: "SliceInt64", tag: `item:"1-3" item-msg:"1-3"`, msg: "1-3", values: []string{"1", "2", "3", "4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceFloat32", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceFloat32", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1", "2", "3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceFloat32", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceFloat32", tag: `item:"1-3" item-msg:"1-3"`, msg: "1-3", values: []string{"1", "2", "3", "4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceFloat64", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceFloat64", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1", "2", "3"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceFloat64", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceFloat64", tag: `item:"1-3" item-msg:"1-3"`, msg: "1-3", values: []string{"1", "2", "3", "4"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "SliceBool", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceBool", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{"1", "1", "1"}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceBool", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{}},
	{line: getFileLine(), field: "SliceBool", tag: `item:"1-3" item-msg:"1-3"`, msg: "1-3", values: []string{"1", "1", "1", "1"}, files: []multipart.FileHeader{}},

	{line: getFileLine(), field: "File", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{}}, // not required
	{line: getFileLine(), field: "File", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{{Filename: "1"}}},
	{line: getFileLine(), field: "File", tag: `item:"1-3" item-msg:"1-3"`, msg: "", values: []string{}, files: []multipart.FileHeader{{Filename: "1"}, {Filename: "2"}, {Filename: "3"}}},
	{line: getFileLine(), field: "File", tag: `item:"1-3" item-msg:"1-3"`, msg: "1-3", values: []string{}, files: []multipart.FileHeader{{Filename: "1"}, {Filename: "2"}, {Filename: "3"}, {Filename: "4"}}},

	// test file size
	{line: getFileLine(), field: "File", tag: `size:"1-3" size-msg:"1-3KB"`, msg: "", values: []string{}, files: []multipart.FileHeader{{Size: 1 << 10}}},
	{line: getFileLine(), field: "File", tag: `size:"1-3" size-msg:"1-3KB"`, msg: "", values: []string{}, files: []multipart.FileHeader{{Size: 3 << 10}}},
	{line: getFileLine(), field: "File", tag: `size-KB:"1-3" size-msg:"1-3KB"`, msg: "1-3KB", values: []string{}, files: []multipart.FileHeader{{Size: 1 << 9}}},
	{line: getFileLine(), field: "File", tag: `size-KB:"1-3" size-msg:"1-3KB"`, msg: "1-3KB", values: []string{}, files: []multipart.FileHeader{{Size: (3 << 10) + 1}}},

	// test file extension
	{line: getFileLine(), field: "File", tag: `exts:".jpg .png .gif" exts-msg:"msg"`, msg: "", values: []string{}, files: []multipart.FileHeader{{Filename: "1.jpg"}, {Filename: "1.png"}, {Filename: "1.gif"}}},
	{line: getFileLine(), field: "File", tag: `exts:".jpg .png .gif" exts-msg:"msg"`, msg: "msg", values: []string{}, files: []multipart.FileHeader{{Filename: "1.jpg"}, {Filename: "1.png"}, {Filename: "1.gif"}, {Filename: "1.txt"}}},

	// test file MIME type
	{line: getFileLine(), field: "File", tag: `mime:"image/png" mime-msg:"msg"`, msg: "", values: []string{}, files: []multipart.FileHeader{
		{Header: textproto.MIMEHeader{"Content-Type": []string{"image/png"}}},
	}},
	{line: getFileLine(), field: "File", tag: `mime:"image/png" mime-msg:"msg"`, msg: "msg", values: []string{}, files: []multipart.FileHeader{
		{Header: textproto.MIMEHeader{"Content-Type": []string{"image/png"}}},
		{Header: textproto.MIMEHeader{"Content-Type": []string{"image/jpeg"}}},
	}},
	{line: getFileLine(), field: "File", tag: `mime:"image" mime-msg:"msg"`, msg: "", values: []string{}, files: []multipart.FileHeader{
		{Header: textproto.MIMEHeader{"Content-Type": []string{"image/png"}}},
		{Header: textproto.MIMEHeader{"Content-Type": []string{"image/jpeg"}}},
	}},
}

// 测试验证器
func TestValidator(t *testing.T) {
	for _, tag := range tags {
		p := &param{tag: &tag}
		parser, err := NewSchema(p, initValidator(p), nil)
		if err != nil {
			t.Error(err)
		} else {

			// 模拟数据
			form := multipart.Form{Value: make(map[string][]string), File: make(map[string][]*multipart.FileHeader)}
			key := tag.field
			if tag.key != "" {
				key = tag.key
			}
			form.Value[key] = tag.values
			for _, file := range tag.files {
				form.File[key] = append(form.File[key], &file)
			}

			// 解析并验证
			newParam := &param{}
			err := parser.Parse(uintptr(unsafe.Pointer(newParam)), &form)
			if tag.msg != "" {
				if cerr, ok := err.(*ValidatorErr); !ok {
					fmt.Printf("%s need: *ConditionErr, got: %v\n", tag.line, err)
				} else {
					if cerr.Message != tag.msg {
						fmt.Printf("%s need:%s, got:%s\n", tag.line, tag.msg, cerr.Message)
						t.Fail()
					}
				}
			} else {
				if err != nil {
					fmt.Printf("need: nil, got: %s", err)
				}
			}
		}
	}
}

func getFileLine() string {
	_, file, line, _ := runtime.Caller(1)
	return fmt.Sprintf("%s:%d", file, line)
}
