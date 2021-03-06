// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package ticker

import (
	"github.com/orivil/morgine/utils/timer"
	"time"
)

var section = timer.NewSectionProvider(5 * time.Second)

func Now() *time.Time {
	return section.Section()
}
