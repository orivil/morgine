// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package param

import (
	"reflect"
	"time"
)

type Kind int

func (k Kind) MarshalJSON() ([]byte, error) {
	return []byte(`"` + k.String() + `"`), nil
}

func (k *Kind) UnmarshalJSON(data []byte) error {
	*k = StrToKind(string(data))
	return nil
}

// supported data type
const (
	Invalid Kind = iota
	Bool
	Int
	Int32
	Int64
	Float32
	Float64
	String
	File         // implement 'param.FileHandler'
	TimePtr      // *time.Time
	SliceString  // []string
	SliceInt     // []int
	SliceInt32   // []int32
	SliceInt64   // []int64
	SliceFloat32 // []float32
	SliceFloat64 // []float64
	SliceBool    // []bool
)

var FieldTypes = map[Kind]string{
	String:       "string",
	Int:          "int",
	Int32:        "int32",
	Int64:        "int64",
	Float32:      "float32",
	Float64:      "float64",
	Bool:         "bool",
	File:         "file",
	TimePtr:      "time",
	SliceString:  "[]string",
	SliceInt:     "[]int",
	SliceInt32:   "[]int32",
	SliceInt64:   "[]int64",
	SliceFloat32: "[]float32",
	SliceFloat64: "[]float64",
	SliceBool:    "[]bool",
	Invalid:      "invalid",
}

func (k Kind) String() string {
	return FieldTypes[k]
}

func StrToKind(str string) Kind {
	for key, value := range FieldTypes {
		if value == str {
			return key
		}
	}
	return Invalid
}

func fieldKind(field reflect.StructField) Kind {
	switch field.Type.Kind() {
	case reflect.String:
		return String
	case reflect.Int:
		return Int
	case reflect.Int32:
		return Int32
	case reflect.Int64:
		return Int64
	case reflect.Float32:
		return Float32
	case reflect.Float64:
		return Float64
	case reflect.Bool:
		return Bool
	case reflect.Ptr:
		var transformerType = reflect.TypeOf(new(time.Time))
		if field.Type.ConvertibleTo(transformerType) {
			return TimePtr
		} else {
			return Invalid
		}
	case reflect.Func:
		var fileHandlerType = reflect.TypeOf(new(FileHandler)).Elem()
		if field.Type.ConvertibleTo(fileHandlerType) {
			return File
		} else {
			return Invalid
		}
	case reflect.Slice:
		element := field.Type.Elem()
		switch element.Kind() {
		case reflect.String:
			return SliceString
		case reflect.Int:
			return SliceInt
		case reflect.Int32:
			return SliceInt32
		case reflect.Int64:
			return SliceInt64
		case reflect.Float32:
			return SliceFloat32
		case reflect.Float64:
			return SliceFloat64
		case reflect.Bool:
			return SliceBool
		}
	}
	return Invalid
}