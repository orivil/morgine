// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import (
	"encoding/json"
	"sync"
)

const (
	// 可根据需要自定义状态码
	StatusSuccess StatusCode = 2000
	StatusNotFound StatusCode = 2404
	StatusUnauthorized StatusCode = 2401
)

// 注册状态码
func init() {
	defaultNamespace := StatusCodes.Namespace(DefaultStatusNamespace)
	defaultNamespace.InitStatus(StatusSuccess,"Success")
	defaultNamespace.InitStatus(StatusNotFound, "NotFound")
	defaultNamespace.InitStatus(StatusUnauthorized,"Unauthorized")
}

const DefaultStatusNamespace = ""

type StatusCode int

var StatusCodes = NewStatusCodes()

// 带状态的数据类型，便于前端判断数据处理方式
type StatusData struct {
	Code StatusCode `json:"code" xml:"code"`
	Data interface{} `json:"data" xml:"data"`
}

type StatusTexts map[StatusCode][]string

type statusCodes struct {
	spaces map[string]*namespaceStatusCodes
	mu sync.Mutex
}

func (sc *statusCodes) StatusTexts() StatusTexts {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	var statusTexts = make(StatusTexts, 10)
	for _, space := range sc.spaces {
		func() {
			space.mu.Lock()
			defer space.mu.Unlock()
			for code, text := range space.Codes {
				if space.Namespace != "" {
					text = space.Namespace + "-" + text
				}
				statusTexts[code] = append(statusTexts[code], text)
			}
		}()
	}
	return statusTexts
}

func (sc *statusCodes) MarshalJSON() ([]byte, error) {
	sts := sc.StatusTexts()
	return json.Marshal(sts)
}

func NewStatusCodes() *statusCodes {
	return &statusCodes{
		spaces: map[string]*namespaceStatusCodes{},
	}
}

func (sc *statusCodes) StatusText(namespace string, code StatusCode) string {
	space := sc.Namespace(namespace)
	return space.StatusText(code)
}

func (sc *statusCodes) InitStatus(namespace string, code StatusCode, codeText string) {
	space := sc.Namespace(namespace)
	space.InitStatus(code, codeText)
}

func (sc *statusCodes) Namespace(namespace string) *namespaceStatusCodes {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	space, ok := sc.spaces[namespace]
	if !ok {
		space = &namespaceStatusCodes {
			Namespace: namespace,
			Codes:     make(map[StatusCode]string, 5),
		}
		sc.spaces[namespace] = space
	}
	return space
}

type namespaceStatusCodes struct {
	Namespace string
	Codes map[StatusCode]string
	mu sync.Mutex
}

func (s *namespaceStatusCodes) StatusText(code StatusCode) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.Codes[code]
}

func (s *namespaceStatusCodes) InitStatus(code StatusCode, codeText string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Codes[code]; ok {
		panic("status code already exists")
	} else {
		s.Codes[code] = codeText
	}
}
