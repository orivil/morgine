// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package utils

import (
	"math/rand"
	"time"
)

// 随机生成一个 min - max 之间的数字
// Example:
//  100000, 999999 => 658459
func RandomInt(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min+1) + min
}
