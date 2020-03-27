// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package limiter

import (
	"golang.org/x/time/rate"
	"math"
	"testing"
	"time"
)

func TestVisitorContainer_Allow(t *testing.T) {
	limitPerMinute := 20  // 限制在1秒钟20次
	totalSecond := 5 // 测试时间 5 秒

	totalAllow := limitPerMinute * totalSecond // 总共被允许的次数

	vis := NewVisitorContainer(func() *rate.Limiter {
		return rate.NewLimiter(rate.Limit(limitPerMinute), 1)
	})
	ip := "192.168.0.1"
	ticker := time.NewTicker(1 * time.Millisecond) // 模拟 1 秒请求 1000 次, 达到限制极限，缩小误差
	allowed := 0
	var endAt *time.Time
	for n := range ticker.C {
		now := n
		if vis.Allow(ip) {
			allowed++
		}
		if endAt == nil {
			end := now.Add(time.Duration(totalSecond) * time.Second)
			endAt = &end
		} else {
			if endAt.Before(now) {
				ticker.Stop()
				break
			}
		}
	}

	var scope float64 = 2  // 误差范围
	if math.Abs(float64(totalAllow - allowed)) >= scope {
		t.Errorf("need ablut %d +/- %f, got %d", totalAllow, scope, allowed)
	}
}
