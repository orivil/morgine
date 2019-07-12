// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package param

import "mime/multipart"

// FileHandler 定义上传文件处理函数, field 是上传文件的字段名, header 为上传文件的信息
type FileHandler func(field string, header *multipart.FileHeader) error
