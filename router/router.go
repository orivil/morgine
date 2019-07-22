// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package router

import (
	"regexp"
	"sort"
	"strings"
	"sync"
)

type Router struct {
	entries map[string]Nodes
	nodes   Nodes
	mu      sync.Mutex
}

func NewRouter() *Router {
	return &Router{entries: make(map[string]Nodes, 8)}
}

func (r *Router) Add(method, route string, action interface{}) error {
	return r.entries[method].add(method, route, action)
}

func (r *Router) Match(method, path string) (action interface{}) {
	if es, ok := r.entries[method]; ok {
		return es.match(path)
	}
	return nil
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

func (es Nodes) Len() int {
	return len(es)
}

func (es Nodes) Less(i, j int) bool {
	return strings.Count(es[i].prefix, "/") < strings.Count(es[j].prefix, "/")
}

func (es Nodes) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

func (es *Nodes) add(method, route string, action interface{}) error {
	prefix, pattern := initRoute(route)
	matcher, err := regexp.Compile(pattern)
	if err != nil {
		return err
	}
	*es = append(*es, &Node{
		prefix:  prefix,
		Route:   route,
		matcher: matcher,
		Action:  action,
		Method:  method,
	})
	sort.Sort(es)
	return nil
}

func (es Nodes) match(path string) (action interface{}) {
	for _, ety := range es {
		if action = ety.match(path); action != nil {
			return action
		}
	}
	return nil
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

func initRoute(route string) (prefix, pattern string) {
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
		if strings.HasSuffix(pattern, "/") {
			pattern = strings.TrimSuffix(pattern, "/")
		} else {
			pattern += "$"
		}
	}
	return
}
