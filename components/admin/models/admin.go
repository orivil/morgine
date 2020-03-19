// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package models

import "time"

type Admin struct {
	ID int
	Account string
	Password string
	CreatedAt *time.Time
}
