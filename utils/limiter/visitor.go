// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package limiter

import (
	"golang.org/x/time/rate"
	"sync"
)

type RateLimiterProvider func() *rate.Limiter

type visitor struct {
	session string
	limiter *rate.Limiter
}

type VisitorContainer struct {
	visitors map[string]*visitor
	limiter RateLimiterProvider
	mu sync.Mutex
}

func NewVisitorContainer(limiter RateLimiterProvider) *VisitorContainer {
	return &VisitorContainer{
		visitors: make(map[string]*visitor, 100),
		limiter: limiter,
	}
}

// session 用于保证用户唯一性，可传入 IP 地址，用户 ID 等，或者使用两个 Container 同时检测 IP 及 ID
func (vc *VisitorContainer) Allow(session string) bool {
	vc.mu.Lock()
	defer vc.mu.Unlock()
	vis := vc.visitors[session]
	if vis == nil {
		vis = &visitor {
			session:          session,
			limiter:      vc.limiter(),
		}
		vc.visitors[session] = vis
		return true
	} else {
		return vis.limiter.Allow()
	}
}
