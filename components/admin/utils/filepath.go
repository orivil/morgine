// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package utils

import (
	"os"
	"path/filepath"
	"strings"
)

type Dir struct {
	Path string
	Name string
	Subs Dirs
}

func (d *Dir) Trim(dir string)  {
	dir = CleanDir(dir)
	d.Path = strings.TrimPrefix(d.Path, dir)
	for _, sub := range d.Subs {
		sub.Trim(dir)
	}
}

func CleanDir(dir string) string {
	return strings.Trim(filepath.ToSlash(dir), "/")
}

type Dirs []*Dir

func initDirs(base string, paths []string) *Dir {
	base = CleanDir(base)
	for i, path := range paths {
		paths[i] = CleanDir(path)
	}
	var dirs Dirs
	for _, path := range paths {
		dirs = append(dirs, &Dir {
			Path: path,
			Name: filepath.Base(path),
		})
	}
	root := &Dir {
		Path: base,
		Name: filepath.Base(base),
	}
	var isSubDir = func(pDir, subDir string) bool {
		if strings.Count(pDir, "/") + 1 == strings.Count(subDir, "/") &&
			strings.HasPrefix(subDir, pDir) {
			return true
		} else {
			return false
		}
	}
	for _, dir := range dirs {
		if isSubDir(root.Path, dir.Path) {
			root.Subs = append(root.Subs, dir)
		} else {
			for _, child := range dirs {
				if isSubDir(dir.Path, child.Path) {
					dir.Subs = append(dir.Subs, child)
				}
			}
		}
	}
	return root
}

// 遍历 dir 目录，获得所有子目录，使用方式参考 example_test.go
func WalkDirs(dir string) (*Dir, error) {
	var paths []string
	var base = CleanDir(dir)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			path = CleanDir(path)
			if base != path {
				paths = append(paths, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return initDirs(base, paths), nil
}