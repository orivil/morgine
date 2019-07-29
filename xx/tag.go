// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"encoding/json"
	"unsafe"
)

type TagName *string

func NewTagName(name string) TagName {
	return &name
}

type ApiTags []*ApiTag

type ApiTag struct {
	Name TagName
	Desc string
	Subs ApiTags
}

func (at *ApiTag) MarshalJSON() ([]byte, error) {
	res := &struct {
		ID   uintptr
		Name TagName
		Desc string
		Subs ApiTags
	}{
		ID:   uintptr(unsafe.Pointer(at.Name)),
		Name: at.Name,
		Desc: at.Desc,
		Subs: at.Subs,
	}
	return json.Marshal(res)
}
