// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package param

import (
	"fmt"
	"mime/multipart"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// 默认时间解析模板
const DefaultTimeLayout = "2006-01-02T15:04:05"

// 时间模板标签名, 设置此标签之后时间按此标签模板解析
const TimeLayoutTag = "time-layout"

type Parser interface {
	Unmarshal(data *multipart.Form, schema interface{}) error
}

type Schema struct {
	Pkg        string
	Name       string
	Fields     []*Field
	setters    []setter
	storeFuncs []setter
}

type Field struct {
	// 字段名称
	Name string

	// 字段描述
	Desc string

	// 默认值
	Value interface{}

	// 字段类型
	Kind Kind
}

type setter func(begin uintptr, form *multipart.Form) error

func NewSchema(v interface{}, validator *Validator, filter *Filter) (*Schema, error) {
	t := reflect.TypeOf(v)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("need struct, got %v", t)
	}
	ptr := reflect.ValueOf(v).Pointer()
	//rt := &Type{}
	//defaultValues := make(map[string][]string, 1)
	schema := &Schema{}
	fields := structFields(t, 0)
	for _, field := range fields {
		if !isFieldIgnore(field) {
			offset := field.Offset
			// 过滤字段
			if filter != nil {
				if b, ok := filter.except[offset]; ok && b {
					continue
				}
				lnOnly := len(filter.only)
				if lnOnly > 0 && !filter.only[offset] {
					continue
				}
			}

			kind := fieldKind(field)
			if kind == Invalid {
				return nil, fmt.Errorf("field [%s] kind is invalid", field.Name)
			}
			f := &Field{
				Name:  fieldName(field),
				Desc:  fieldDesc(field),
				Value: fieldDefaultValue(kind, ptr, field.Offset),
				Kind:  kind,
			}

			var cdt *condition
			// 添加接口条件
			if validator != nil {
				if c, ok := validator.conditions[field.Offset]; ok {
					cdt = c
				}
			}
			// 接口条件不存在, 则尝试添加标签条件
			if cdt == nil {
				if isCondition(field.Tag) {
					c := &condition{}
					err := c.Syntax(string(field.Tag))
					if err != nil {
						// 语法错误
						return nil, fmt.Errorf("%T.%s error:%s", v, field.Name, err)
					} else {
						cdt = c
					}
				}
			}
			switch kind {
			case File:
			default:

			}
			schema.Fields = append(schema.Fields, f)
		}
	}
	return schema, nil
}

// 获取 struct 的字段偏移量及字段类型, 包括嵌套的字段. 只支持嵌套的 struct, 不支持 struct ptr
func structFields(structural reflect.Type, fieldsOffset uintptr) (fields []reflect.StructField) {
	fieldNum := structural.NumField()
	for idx := 0; idx < fieldNum; idx++ {
		field := structural.Field(idx)
		// 嵌套 struct 需要加上父字段的偏移量
		field.Offset += fieldsOffset
		kind := field.Type.Kind()
		if kind == reflect.Struct { // 只支持 struct, 不支持 struct ptr
			subFields := structFields(field.Type, field.Offset)
			fields = append(fields, subFields...)
		} else {
			fields = append(fields, field)
		}
	}
	return
}

func isFieldIgnore(field reflect.StructField) bool {
	return field.Tag.Get("param") == "-"
}

func fieldName(field reflect.StructField) string {
	name := field.Tag.Get("param")
	if name != "" {
		return name
	}
	return field.Name
}

func fieldDesc(field reflect.StructField) string {
	return field.Tag.Get("desc")
}

func fieldDefaultValue(kind Kind, ptr, offset uintptr) interface{} {
	switch kind {
	case String:
		return (*string)(unsafe.Pointer(ptr + offset))
	case Int:
		return (*int)(unsafe.Pointer(ptr + offset))
	case Int32:
		return (*int32)(unsafe.Pointer(ptr + offset))
	case Int64:
		return (*int64)(unsafe.Pointer(ptr + offset))
	case Float32:
		return (*float32)(unsafe.Pointer(ptr + offset))
	case Float64:
		return (*float64)(unsafe.Pointer(ptr + offset))
	case Bool:
		return (*bool)(unsafe.Pointer(ptr + offset))
	case File:
		handler := (*FileHandler)(unsafe.Pointer(ptr + offset))
		if handler != nil {
			return handler
		}
		var f FileHandler = func(field string, header *multipart.FileHeader) error {
			return fmt.Errorf("parameter field [%s], the file handler is nil", field)
		}
		return f
	case TimePtr:
		def := *(**time.Time)(unsafe.Pointer(ptr + offset))
		if def != nil {
			return def
		}
		return &time.Time{}
	case SliceString:
		return (*[]string)(unsafe.Pointer(ptr + offset))
	case SliceInt:
		return (*[]int)(unsafe.Pointer(ptr + offset))
	case SliceInt32:
		return (*[]int32)(unsafe.Pointer(ptr + offset))
	case SliceInt64:
		return (*[]int64)(unsafe.Pointer(ptr + offset))
	case SliceFloat32:
		return (*[]float32)(unsafe.Pointer(ptr + offset))
	case SliceFloat64:
		return (*[]float64)(unsafe.Pointer(ptr + offset))
	case SliceBool:
		return (*[]bool)(unsafe.Pointer(ptr + offset))
	}
}

func getSetter(field reflect.StructField, param string, kind Kind, ptr, offset uintptr, dvalue interface{}, cdt *condition) setter {
	switch kind {
	case String:
		return newStringSetter(param, offset, dvalue.(string), cdt)
	case Int:
		return newIntSetter(param, offset, dvalue.(int), cdt)
	case Int32:
		return newInt32Setter(param, offset, dvalue.(int32), cdt)
	case Int64:
		return newInt64Setter(param, offset, dvalue.(int64), cdt)
	case Float32:
		return newFloat32Setter(param, offset, dvalue.(float32), cdt)
	case Float64:
		return newFloat64Setter(param, offset, dvalue.(float64), cdt)
	case Bool:
		return newBoolSetter(param, offset, dvalue.(bool))
	case File:
		return newFileSetter(param, offset, dvalue.(FileHandler), cdt)
	case TimePtr:
		layout := field.Tag.Get(TimeLayoutTag)
		if layout == "" {
			layout = DefaultTimeLayout
		}
		return newTimeSetter(param, layout, offset, dvalue.(*time.Time), cdt)
	case SliceString:
		setter = newSliceSetter(SliceString, param, offset, cdt)
		info.Value = (*[]string)(unsafe.Pointer(ptr + offset))
	case SliceInt:
		setter = newSliceSetter(SliceInt, param, offset, cdt)
		info.Value = (*[]int)(unsafe.Pointer(ptr + offset))
	case SliceInt32:
		setter = newSliceSetter(SliceInt32, param, offset, cdt)
		info.Value = (*[]int32)(unsafe.Pointer(ptr + offset))
	case SliceInt64:
		setter = newSliceSetter(SliceInt64, param, offset, cdt)
		info.Value = (*[]int64)(unsafe.Pointer(ptr + offset))
	case SliceFloat32:
		setter = newSliceSetter(SliceFloat32, param, offset, cdt)
		info.Value = (*[]float32)(unsafe.Pointer(ptr + offset))
	case SliceFloat64:
		setter = newSliceSetter(SliceFloat64, param, offset, cdt)
		info.Value = (*[]float64)(unsafe.Pointer(ptr + offset))
	case SliceBool:
		setter = newSliceSetter(SliceBool, param, offset, cdt)
		info.Value = (*[]bool)(unsafe.Pointer(ptr + offset))
	}
}

func getStoreFunc(field reflect.StructField) setter {

}

func newStringSetter(param string, offset uintptr, dvalue string, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		value := url.Values(form.Value).Get(param)
		if value == "" {
			value = dvalue
		}
		*(*string)(unsafe.Pointer(begin + offset)) = value
		if cdt != nil {
			err = cdt.validStr(param, value)
			if err == nil {
				err = cdt.validEnum(param, form.Value[param])
			}
			if err != nil {
				return err
			}
		}
		return
	}
}

func newIntSetter(param string, offset uintptr, dvalue int, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		valueStr := url.Values(form.Value).Get(param)
		var value int
		if valueStr == "" {
			if dvalue != 0 {
				value = dvalue
			}
		} else {
			value, err = strconv.Atoi(valueStr)
			if err != nil {
				return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
			}
		}
		*(*int)(unsafe.Pointer(begin + offset)) = value

		if cdt != nil {
			err = cdt.validNum(param, valueStr, float64(value))
			if err == nil {
				err = cdt.validEnum(param, form.Value[param])
			}
			if err != nil {
				return err
			}
		}
		return
	}
}

func newInt32Setter(param string, offset uintptr, dvalue int32, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		valueStr := url.Values(form.Value).Get(param)
		var value int32
		if valueStr == "" {
			if dvalue != 0 {
				value = dvalue
			}
		} else {
			i64, err := strconv.ParseInt(valueStr, 10, 32)
			if err != nil {
				return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
			}
			value = int32(i64)
		}
		*(*int32)(unsafe.Pointer(begin + offset)) = value

		if cdt != nil {
			err = cdt.validNum(param, valueStr, float64(value))
			if err == nil {
				err = cdt.validEnum(param, form.Value[param])
			}
			if err != nil {
				return err
			}
		}
		return
	}
}

func newInt64Setter(param string, offset uintptr, dvalue int64, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		valueStr := url.Values(form.Value).Get(param)
		var value int64
		if valueStr == "" {
			if dvalue != 0 {
				value = dvalue
			}
		} else {
			value, err = strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
			}
		}
		*(*int64)(unsafe.Pointer(begin + offset)) = value

		if cdt != nil {
			err = cdt.validNum(param, valueStr, float64(value))
			if err == nil {
				err = cdt.validEnum(param, form.Value[param])
			}
			if err != nil {
				return err
			}
		}
		return
	}
}

func newFloat64Setter(param string, offset uintptr, dvalue float64, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		valueStr := url.Values(form.Value).Get(param)
		var value float64
		if valueStr == "" {
			if dvalue != 0 {
				value = dvalue
			}
		} else {
			value, err = strconv.ParseFloat(valueStr, 64)
			if err != nil {
				return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
			}
		}
		*(*float64)(unsafe.Pointer(begin + offset)) = value

		if cdt != nil {
			err = cdt.validNum(param, valueStr, value)
			if err == nil {
				err = cdt.validEnum(param, form.Value[param])
			}
			if err != nil {
				return err
			}
		}
		return
	}
}

func newFloat32Setter(param string, offset uintptr, dvalue float32, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		valueStr := url.Values(form.Value).Get(param)
		var value float32
		if valueStr == "" {
			if dvalue != 0 {
				value = dvalue
			}
		} else {
			f64, err := strconv.ParseFloat(valueStr, 32)
			if err != nil {
				return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
			}
			value = float32(f64)
		}
		*(*float32)(unsafe.Pointer(begin + offset)) = value

		if cdt != nil {
			err = cdt.validNum(param, valueStr, float64(value))
			if err == nil {
				err = cdt.validEnum(param, form.Value[param])
			}
			if err != nil {
				return err
			}
		}
		return
	}
}

func newBoolSetter(param string, offset uintptr, dvalue bool) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		valueStr := url.Values(form.Value).Get(param)
		value := valueStr != "" && valueStr != "false"
		*(*bool)(unsafe.Pointer(begin + offset)) = value
		return
	}
}

func newFileStorage(param string, offset uintptr) storage {
	return func(begin uintptr, form *multipart.Form) (err error) {
		handler := *(*FileHandler)(unsafe.Pointer(begin + offset))
		if handler != nil {
			for _, header := range form.File[param] {
				err = handler(param, header)
				if err != nil {
					return err
				}
			}
		}
		return nil
	}
}

func newFileSetter(param string, offset uintptr, dvalue FileHandler, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		if cdt != nil {
			return cdt.validFile(param, form)
		}
		handler := *(*FileHandler)(unsafe.Pointer(begin + offset))
		if handler == nil {
			handler = dvalue
		}
		for _, header := range form.File[param] {
			err = handler(param, header)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

func newTimeSetter(param, layout string, offset uintptr, dvalue *time.Time, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		if cdt != nil {
			err = cdt.validTime(param, form)
			if err != nil {
				return err
			}
		}
		text := url.Values(form.Value).Get(param)
		var t time.Time
		if len(text) > 0 {
			t, err = time.Parse(layout, text)
			if err != nil {
				return
			}
		} else {
			t = time.Unix(0, dvalue.UnixNano())
		}
		*(**time.Time)(unsafe.Pointer(begin + offset)) = &t
		return
	}
}

func newSliceStringSetter(param string, offset uintptr, dvalue []string, cdt *condition) {

}

func newSliceSetter(kind Kind, param string, offset uintptr, cdt *condition) setter {
	return func(begin uintptr, form *multipart.Form) (err error) {
		values := getSliceValues(param, form)
		if len(values) > 0 {
			switch kind {
			case SliceString:
				*(*[]string)(unsafe.Pointer(begin + offset)) = values
			case SliceInt:
				var ints = make([]int, len(values))
				for idx, value := range values {
					ints[idx], err = strconv.Atoi(value)
					if err != nil {
						return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
					}
				}
				*(*[]int)(unsafe.Pointer(begin + offset)) = ints
			case SliceInt32:
				var ints = make([]int32, len(values))
				for idx, value := range values {
					i32, err := strconv.ParseInt(value, 10, 32)
					if err != nil {
						return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
					}
					ints[idx] = int32(i32)
				}
				*(*[]int32)(unsafe.Pointer(begin + offset)) = ints
			case SliceInt64:
				var ints = make([]int64, len(values))
				for idx, value := range values {
					i64, err := strconv.ParseInt(value, 10, 64)
					if err != nil {
						return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
					}
					ints[idx] = int64(i64)
				}
				*(*[]int64)(unsafe.Pointer(begin + offset)) = ints
			case SliceFloat32:
				var floats = make([]float32, len(values))
				for idx, value := range values {
					f32, err := strconv.ParseFloat(value, 32)
					if err != nil {
						return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
					}
					floats[idx] = float32(f32)
				}
				*(*[]float32)(unsafe.Pointer(begin + offset)) = floats
			case SliceFloat64:
				var floats = make([]float64, len(values))
				for idx, value := range values {
					f64, err := strconv.ParseFloat(value, 64)
					if err != nil {
						return &ValidatorErr{Field: param, Kind: ConditionInvalidNumber, Message: err.Error()}
					}
					floats[idx] = f64
				}
				*(*[]float64)(unsafe.Pointer(begin + offset)) = floats
			case SliceBool:
				var bools = make([]bool, len(values))
				for idx, value := range values {
					bools[idx] = value != "" && value != "false"
				}
				*(*[]bool)(unsafe.Pointer(begin + offset)) = bools
			}
		}
		if cdt != nil {
			err = cdt.validItem(param, len(values))
			if err == nil {
				err = cdt.validEnum(param, values)
			}
			if err != nil {
				return err
			}
		}
		return
	}
}

func getSliceValues(param string, form *multipart.Form) (vs []string) {
	// normal type: foo[]=1&foo[]=2
	vs = form.Value[param+"[]"]

	if len(vs) == 0 {
		// traditional type: foo=1&foo=2
		vs = form.Value[param]
	}
	// other type: foo=1,2
	if len(vs) == 1 && strings.Contains(vs[0], ",") {
		return strings.Split(vs[0], ",")
	} else {
		return vs
	}
}
