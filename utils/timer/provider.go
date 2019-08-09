// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package timer

import (
	"sync"
	"time"
)

// 时间段提供器, 用于某些高并发需求, 避免过多调用系统时间
type SectionProvider struct {
	section *time.Time
	locker  sync.RWMutex
	close   chan struct{}
}

func NewSectionProvider(section time.Duration) *SectionProvider {
	now := time.Now()
	p := &SectionProvider{section: &now, close: make(chan struct{})}
	go func() {
		ticker := time.NewTicker(section)
		for {
			select {
			case now := <-ticker.C:
				p.locker.Lock()
				p.section = &now
				p.locker.Unlock()
			case <-p.close:
				ticker.Stop()
				return
			}
		}
	}()
	return p
}

func (p *SectionProvider) Close() {
	p.close <- struct{}{}
}

// 提供一个时间段
func (p *SectionProvider) Section() *time.Time {
	p.locker.RLock()
	defer p.locker.RUnlock()
	return p.section
}
