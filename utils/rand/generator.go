// Copyright 2018 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package rand

import (
	"crypto/rand"
)

// 允许的字符集
var str = []byte("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 创建指定位数的随机字符数组
func NewRandByte(bit int) []byte {
	return NewSrcRandByte(str, bit)
}

// 从src中随机创建bit位的随机字符数组
func NewSrcRandByte(src []byte, bit int) []byte {
	bts := make([]byte, bit)
	r := make([]byte, bit)
	ln := len(src)
	rand.Read(bts)
	for idx, b := range bts {
		r[idx] = src[int(b)%ln]
	}
	return r
}
