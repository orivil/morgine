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

func (tags ApiTags) checkIsSubTag(tag TagName) bool {
	for _, at := range tags {
		if at.Name == tag {
			return true
		}
		if at.Subs != nil {
			exist := at.Subs.checkIsSubTag(tag)
			if exist {
				return true
			}
		}
	}
	return false
}

func (tags ApiTags) checkIsEndTag(tag TagName) bool {
	for _, at := range tags {
		if at.Name == tag {
			if len(at.Subs) == 0 {
				return true
			} else {
				return false
			}
		}
		if at.Subs != nil {
			if at.Subs.checkIsEndTag(tag) {
				return true
			}
		}
	}
	return false
}
