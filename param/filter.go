// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package param

import "reflect"

type Filter struct {
	ptr    uintptr
	except map[uintptr]bool
	only   map[uintptr]bool
}

func (f *Filter) Except(pointer ...interface{}) *Filter {
	for _, p := range pointer {
		ptr := reflect.ValueOf(p).Pointer()
		offset := ptr - f.ptr
		f.except[offset] = true
	}
	return f
}

func (f *Filter) Only(pointer ...interface{}) *Filter {
	for _, p := range pointer {
		ptr := reflect.ValueOf(p).Pointer()
		offset := ptr - f.ptr
		f.only[offset] = true
	}
	return f
}

func NewFilter(v interface{}) *Filter {
	return &Filter{
		ptr:    reflect.ValueOf(v).Pointer(),
		except: make(map[uintptr]bool, 2),
		only:   make(map[uintptr]bool, 2),
	}
}
