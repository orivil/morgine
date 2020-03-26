// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package limiter

import (
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type visitor struct {
	session string
	limiter *rate.Limiter
}

type VisitorContainer struct {
	visitors map[string]*visitor
	mu sync.Mutex
}

func NewVisitorContainer() *VisitorContainer {
	return &VisitorContainer{
		visitors: make(map[string]*visitor, 100),
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
			limiter:      rate.NewLimiter(rate.Every(time.Millisecond * 50), 1), // 1秒钟超过 20 次触发等待
		}
		vc.visitors[session] = vis
		return true
	} else {
		return vis.limiter.Allow()
	}
}
