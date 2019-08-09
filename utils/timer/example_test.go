// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package timer_test

import (
	"fmt"
	"github.com/orivil/morgine/utils/timer"
	"time"
)

func ExampleNewTickerRunner() {
	now := time.Now()
	// yesterday start time
	t := timer.LocalHourTime(-24, now)
	fmt.Println(t)
	// Output: 2019-04-04 00:00:00 +0800 CST

	// today start time
	t = timer.LocalHourTime(0, now)
	fmt.Println(t)
	// Output: 2019-04-05 00:00:00 +0800 CST

	// tomorrow start time
	t = timer.LocalHourTime(24, now)
	fmt.Println(t)
	// Output: 2019-04-06 00:00:00 +0800 CST

	// previous hour start time
	t = timer.MinuteTime(-60, now)
	fmt.Println(t)
	// Output: 2019-04-05 13:00:00 +0800 CST

	// this hour start time
	t = timer.MinuteTime(0, now)
	fmt.Println(t)
	// Output: 2019-04-05 14:00:00 +0800 CST

	// next hour start time
	t = timer.MinuteTime(60, now)
	fmt.Println(t)
	// Output: 2019-04-05 15:00:00 +0800 CST

	// previous minute start time
	t = timer.SecondTime(-60, now)
	fmt.Println(t)
	// Output: 2019-04-05 14:23:00 +0800 CST

	// this minute start time
	t = timer.SecondTime(0, now)
	fmt.Println(t)
	// Output: 2019-04-05 14:24:00 +0800 CST

	// next minute start time
	t = timer.SecondTime(60, now)
	fmt.Println(t)
	// Output: 2019-04-05 14:25:00 +0800 CST

	// 第二天凌晨 0 点开始执行一次, 且启动 ticker
	startTime := timer.LocalHourTime(24, now)

	// 触发周期
	tickerDuration := 24 * time.Hour

	// 加入 callback 的时候是否立即执行一次
	startImmediately := false

	// 新建执行器
	runner := timer.NewTickerRunner(tickerDuration, startTime, startImmediately)

	// 向执行器中加入回调函数
	runner.AddCallback(func(now *time.Time) {
		fmt.Println("now:", now)
	})

	// Output:
	// ...
}
