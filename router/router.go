// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package router

import (
	"container/heap"
	"net/url"
	"regexp"
	"strings"
	"sync"
)

type Values func() url.Values

type Router struct {
	entries map[string]Nodes
	nodes   Nodes
	mu      sync.Mutex
}

func NewRouter() *Router {
	return &Router{entries: make(map[string]Nodes, 8)}
}

func (r *Router) Add(method, route string, action interface{}) error {
	if _, ok := r.entries[method]; !ok {
		r.entries[method] = make(Nodes, 0)
	}
	nodes := r.entries[method]
	err := nodes.add(method, route, action)
	if err != nil {
		return err
	} else {
		r.entries[method] = nodes
		return nil
	}
}

func (r *Router) Match(method, path string) (vs Values, action interface{}) {
	if es, ok := r.entries[method]; ok {
		matcher, act := es.match(path)
		return func() url.Values { return getValues(matcher, path) }, act
	}
	return nil, nil
}

func (r *Router) Nodes() Nodes {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.nodes == nil {
		for _, es := range r.entries {
			r.nodes = append(r.nodes, es...)
		}
	}
	return r.nodes
}

type Nodes []*Node

func (es *Nodes) Push(v interface{}) {
	*es = append(*es, v.(*Node))
}

func (es *Nodes) Pop() (v interface{}) {
	*es, v = (*es)[:es.Len()-1], (*es)[es.Len()-1]
	return
}

func (es Nodes) Len() int {
	return len(es)
}

func (es Nodes) Less(i, j int) bool {
	ci := strings.Count(es[i].prefix, "/")
	cj := strings.Count(es[j].prefix, "/")
	if ci == cj {
		return len(es[i].prefix) > len(es[j].prefix)
	}
	return ci > cj
}

func (es Nodes) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es *Nodes) add(method, route string, action interface{}) error {
	prefix, pattern := InitRoute(route)
	matcher, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	heap.Push(es, &Node{
		prefix:  prefix,
		Route:   route,
		matcher: matcher,
		Action:  action,
		Method:  method,
	})
	return nil
}

func (es Nodes) match(path string) (matcher *regexp.Regexp, action interface{}) {
	for _, ety := range es {
		if action = ety.match(path); action != nil {
			return ety.matcher, action
		}
	}
	return nil, nil
}

type Node struct {
	Method  string
	Route   string
	Action  interface{}
	prefix  string
	matcher *regexp.Regexp
}

func (ety *Node) match(path string) (action interface{}) {
	if !strings.HasPrefix(path, ety.prefix) {
		return nil
	}
	if ety.matcher != nil && ety.matcher.MatchString(path) {
		return ety.Action
	}
	return nil
}

var paramPatternReplacer = regexp.MustCompile("{[^\\/]+?}")

var pathPatternMatcher = regexp.MustCompile("^/[\\w|\\-|\\_|\\/|\\.]*")

// 将路由格式转换成正则匹配格式
func InitRoute(route string) (prefix, pattern string) {
	switch route {
	case "", "/":
		prefix = "/"
		pattern = "^/$"
	default:
		prefix = pathPatternMatcher.FindString(route)
		if prefix != "/" {
			prefix = strings.TrimSuffix(prefix, "/")
		}
		pattern = paramPatternReplacer.ReplaceAllStringFunc(route, func(s string) string {
			s = strings.TrimPrefix(s, "{")
			s = strings.TrimSuffix(s, "}")
			return "(?P<" + s + ">[^\\/^\\.]+)"
		})
		pattern = "^" + pattern
		if !strings.HasSuffix(pattern, "/") {
			pattern += "$"
		}
	}
	return
}

// 获得 path 中的参数
func getValues(matcher *regexp.Regexp, path string) url.Values {
	vs := make(url.Values, 2)
	names := matcher.SubexpNames()
	values := matcher.FindStringSubmatch(path)
	for key, name := range names {
		if name != "" && values[key] != "" {
			vs[name] = append(vs[name], values[key])
		}
	}
	return vs
}
