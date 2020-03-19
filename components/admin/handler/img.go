// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package handler

import (
	"github.com/orivil/morgine/components/admin/env"
	"github.com/orivil/morgine/xx"
	"os"
	"path/filepath"
	"strings"
)



func GetImageDirs(method, route string, cdt *xx.Condition) {
	doc := &xx.Doc {
		Title:     "获得目录列表",
		Desc:      "",
		Params:    nil,
		Responses: nil,
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		var dirs []
		err := filepath.Walk(env.Config.ImgDir, func(path string, info os.FileInfo, err error) error {

		})
		if err != nil {
			ctx.Error(err)
		}
	})
}
