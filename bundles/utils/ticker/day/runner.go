// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package day_ticker

import (
	"github.com/orivil/morgine/utils/timer"
)

const (
	// 0-23
	TickerStartHour = 0 // 每天 4 点开始
)

var Runner = timer.NewDayTicker(TickerStartHour, false)
