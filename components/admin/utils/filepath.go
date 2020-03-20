// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package utils

import (
	"os"
	"path/filepath"
	"strings"
)

// 遍历 dir 目录，获得所有子目录，使用方式参考 example_test.go
func WalkDirs(dir string) (dirs map[string]interface{}, err error) {
	dirs = make(map[string]interface{})
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			path = filepath.ToSlash(path)
			paths := strings.Split(path, "/")
			var walk func(dirs map[string]interface{}, paths []string)
			walk = func(dirs map[string]interface{}, paths []string) {
				if ln := len(paths); ln > 0 {
					fp := paths[0]
					if si, ok := dirs[fp]; ok {
						var sub map[string]interface{}
						if si == nil {
							sub = make(map[string]interface{})
						} else {
							sub = si.(map[string]interface{})
						}
						dirs[fp] = sub
						if ln > 1 {
							walk(sub, paths[1:])
						}
					} else {
						dirs[fp] = make(map[string]interface{})
					}
				}
			}
			walk(dirs, paths)
		}
		return nil
	})
	return dirs, err
}