// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package hour_ticker

import (
	"github.com/orivil/morgine/utils/timer"
	"time"
)

var Runner = timer.NewTickerRunner(1*time.Hour, nil, false)
