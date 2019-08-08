// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"encoding/json"
	"fmt"
	"github.com/orivil/morgine/param"
	"runtime"
	"sort"
	"unsafe"
)

type apiDoc struct {
	Tags    ApiTags
	Middles map[uintptr]*apiMiddle
	Actions map[uintptr][]*apiAction
}

func newApiDoc() *apiDoc {
	return &apiDoc{
		Middles: map[uintptr]*apiMiddle{},
		Actions: map[uintptr][]*apiAction{},
	}
}

func (doc *apiDoc) use(middles ...*Handler) {
	for _, middle := range middles {
		ptr := uintptr(unsafe.Pointer(middle))
		if _, ok := doc.Middles[ptr]; !ok {
			doc.Middles[ptr] = &apiMiddle{
				Name:      middle.Doc.Title,
				Desc:      middle.Doc.Desc,
				Params:    initApiParams(middle.Doc.parser),
				Responses: middle.Doc.Responses,
			}
		}
	}
}

func (doc *apiDoc) handle(depth int, tag TagName, method, route string, d *Doc, middles []*Handler) {
	act := &apiAction{
		Name:        d.Title,
		Desc:        d.Desc,
		Trace:       initTrace(depth + 1),
		Method:      method,
		Route:       route,
		Params:      initApiParams(d.parser),
		ContentType: getActionContentType(d.parser),
		Responses:   d.Responses,
	}
	for _, middle := range middles {
		act.Middles = append(act.Middles, uintptr(unsafe.Pointer(middle)))
	}
	ptr := uintptr(unsafe.Pointer(tag))
	doc.Actions[ptr] = append(doc.Actions[ptr], act)
}

func initTrace(depth int) string {
	_, file, line, _ := runtime.Caller(depth + 1)
	return fmt.Sprintf("%s: %d", file, line)
}

type apiMiddle struct {
	Name      string
	Desc      string
	Params    []*apiParam
	Responses Responses
}

type apiParam struct {
	Type   ParamType
	Fields []*param.Field
}

type apiParams []*apiParam

func (ps apiParams) Len() int {
	return len(ps)
}

func (ps apiParams) Less(i, j int) bool {
	return ps[i].Type < ps[j].Type
}

func (ps apiParams) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func getActionContentType(p *parser) param.EncodeType {
	for _, schema := range p.schemas {
		if schema.EncodeType() == param.FormDataEncodeType {
			return param.FormDataEncodeType
		}
	}
	return param.UrlEncodeType
}

func initApiParams(p *parser) apiParams {
	var params apiParams
	for key, typ := range p.types {
		schema := p.schemas[key]
		p := &apiParam{
			Type:   typ,
			Fields: schema.Fields,
		}
		params = append(params, p)
	}
	sort.Sort(params)
	return params
}

type apiAction struct {
	Name        string
	Desc        string
	Trace       string
	Method      string
	Route       string
	Middles     []uintptr
	Params      []*apiParam
	ContentType param.EncodeType
	Responses   Responses
}

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
