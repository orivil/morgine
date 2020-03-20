// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package xx

import "fmt"

type StatusCode int

func (sc StatusCode) Name() string {
	return StatusCodes[sc]
}

func (sc StatusCode) Error() string {
	return sc.Name()
}

var StatusCodes map[StatusCode]string

func InitStatusCode(code StatusCode, name string) {
	if _, ok := StatusCodes[code]; ok {
		panic(fmt.Errorf("status %d already exist", code))
	} else {
		StatusCodes[code] = name
	}
}

const (
	// 可根据需要自定义状态码
	StatusSuccess StatusCode = 2000
	StatusNotFound StatusCode = 2404
	StatusUnauthorized StatusCode = 2401
	StatusTokenInvalid StatusCode = 2402
	StatusTokenExpired StatusCode = 2403
)

func init() {
	InitStatusCode(StatusSuccess, "Success")
	InitStatusCode(StatusNotFound, "NotFound")
	InitStatusCode(StatusUnauthorized, "Unauthorized")
	InitStatusCode(StatusTokenInvalid, "TokenInvalid")
	InitStatusCode(StatusTokenExpired, "TokenExpired")
}
