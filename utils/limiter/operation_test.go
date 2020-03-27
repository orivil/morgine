// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package limiter

import (
	"sync"
	"testing"
	"time"
)

type waitTime struct {
	now time.Time
	mu sync.Mutex
}

func (w *waitTime) GetNowTime() time.Time {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.now
}

func (w *waitTime) GetWaitTime(failed int) time.Duration {
	switch failed {
	case 1, 2, 3:
		return 0
	case 4:
		return 1 * time.Minute
	default:
		return time.Duration((failed - 4) * 2) * time.Minute
	}
}

func (w *waitTime) setNowTime(now time.Time) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.now = now
}

func TestNewOperationContainer(t *testing.T) {
	waitTimeProvider := &waitTime{}
	opc := NewOperationContainer(waitTimeProvider)
	ip := "127.0.0.1"
	for failed := 1; failed < 100; failed++ {
		now := time.Now() // 统一当前时间
		waitTimeProvider.setNowTime(now)
		opc.Failed(ip)
		if opc.Allow(ip) {
			got := opc.GetWaitTime(ip)
			if nil != got {
				t.Errorf("%d need: %v got: %v\n", failed, nil, got)
			}
			wait := waitTimeProvider.GetWaitTime(failed)
			if wait != 0 {
				t.Errorf("need: %v got: %v\n", 0, got)
			}
		} else {
			waitDuration := waitTimeProvider.GetWaitTime(failed)
			waitAt := opc.GetWaitTime(ip)
			need := now.Add(waitDuration)
			got := *waitAt
			if !got.Equal(need) {
				t.Errorf("need: %v got: %v\n", need, got)
			}
		}
	}
}

