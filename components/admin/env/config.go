// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package env

var Config = struct {
	// 图片顶级目录, 后面不可加 “/”
	ImgDir string
	AuthTokenExpiredDay int
}{
	ImgDir: "images",
	AuthTokenExpiredDay: 30,
}