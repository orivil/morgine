// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package timer

import (
	"github.com/orivil/morgine/log"
	"runtime"
	"sync"
	"time"
)

type Callback func(now *time.Time)

type TickerRunner struct {
	signals          []chan *time.Time
	startImmediately bool
	close            chan struct{}
	mu               sync.RWMutex
}

func (tr *TickerRunner) Close() {
	close(tr.close)
}

func (tr *TickerRunner) AddCallback(call Callback) {
	if tr.startImmediately {
		now := time.Now()
		call(&now)
	}
	signal := make(chan *time.Time, 2)
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				log.Panic.Printf("panic: %v \n%s\n", err, buf)
			}
		}()
		for {
			select {
			case now := <-signal:
				call(now)
			case <-tr.close:
				return
			}
		}
	}()
	tr.mu.Lock()
	defer tr.mu.Unlock()
	tr.signals = append(tr.signals, signal)
}

func NewTickerRunner(heartbeat time.Duration, startTime *time.Time, startImmediately bool) *TickerRunner {
	runner := &TickerRunner{
		startImmediately: startImmediately,
		close:            make(chan struct{}),
	}
	timeTicker := time.NewTicker(heartbeat)
	go func() {
		call := func(now *time.Time) {
			runner.mu.RLock()
			defer runner.mu.RUnlock()
			for _, signal := range runner.signals {
				signal <- now
			}
		}
		if startTime != nil {
			sleepTime := startTime.Unix() - time.Now().Unix()
			if sleepTime > 0 {
				time.Sleep(time.Duration(sleepTime) * time.Second)
				now := time.Now()
				call(&now)
			}
		}
		for {
			select {
			case now := <-timeTicker.C:
				call(&now)
			case <-runner.close:
				return
			}
		}
	}()
	return runner
}

// yesterday start time: -24
// today start time: 0
// tomorrow start time: 24
func LocalHourTime(hour int, now time.Time) *time.Time {
	now = now.Local()
	y, m, d := now.Date()
	today := time.Date(y, m, d, hour, 0, 0, 0, now.Location())
	return &today
}

// previous hour start time: -60
// this hour start time: 0
// next hour start time: 60
func MinuteTime(minute int, now time.Time) *time.Time {
	y, m, d := now.Date()
	today := time.Date(y, m, d, now.Hour(), minute, 0, 0, now.Location())
	return &today
}

// previous minute start time: -60
// this minute start time: 0
// next minute start time: 60
func SecondTime(second int, now time.Time) *time.Time {
	y, m, d := now.Date()
	today := time.Date(y, m, d, now.Hour(), now.Minute(), second, 0, now.Location())
	return &today
}

// NewDayTicker 用于新建一个每天定时触发的触发器.
// startHour 代表开始触发的时间点(小时), 必须在 0-23 之间
// runImmediately 代表是否在加入回调函数时立即执行一次
func NewDayTicker(startHour int, runImmediately bool) *TickerRunner {
	if 0 <= startHour && startHour <= 23 {
		now := time.Now()
		nowHour := now.Hour()
		if startHour <= nowHour {
			startHour += 24 // 如果当前时间超过设定的时间, 则第二天开始执行
		}
		return NewTickerRunner(time.Duration(24)*time.Hour, LocalHourTime(startHour, now), runImmediately)
	} else {
		panic("start hour should between 0-23")
	}
}
