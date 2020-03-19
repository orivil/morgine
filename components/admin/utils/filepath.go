// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Dir struct {
	Subs Dirs
	Path string
}

type Dirs []*Dir

func (ds Dirs) Add(path string) {
	path = filepath.ToSlash(path)
	sp := strings.Split(path, "/")
	ds.append(sp)
}

func (ds *Dirs) append(splitPath []string) {
	//fmt.Println(splitPath)
	if ln := len(splitPath); ln > 0 {
		firstPart := splitPath[0]
		for _, d := range *ds {
			if d.Path == firstPart {
				d.Path = firstPart
				if ln > 1 {
					d.Subs.append(splitPath[1:])
				}
				return
			}
		}
		*ds = append(*ds, &Dir{Path:firstPart})
	}
}

func WalkDirs(path string) (dirs interface{}, err error) {
	var list []string
	var dirs map[string]interface{}
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			path = filepath.ToSlash(path)
			paths := filepath.SplitList(path)
			var walk func(dirs map[string]interface{}, paths []string)
			walk = func(dirs map[string]interface{}, paths []string) {
				if len(paths) > 0 {
					fp := paths[0]
					if si, ok := dirs[fp]; ok {
						var sub map[string]interface{}
						if si == nil {
							sub = make(map[string]interface{})
						} else {
							sub = si.(map[string]interface{})
						}
						dirs[]
						walk(sub, paths[1:])
					} else {
						dirs[fp] = struct {}{}
					}
				}
			}
			for lidx := len(paths) - 1; lidx >=0; lidx -- {
				p := paths[lidx]

			}
			//list = append(list, p)
			//for lastIdx := len(list) - 1; lastIdx >= 0; lastIdx-- {
			//	lp := list[lastIdx]
			//	if isChildDir(lp, )
			//}
		}
		return nil
	})
	return dirs, err
}

func isChildDir(parent, child string) bool {
	return strings.Count(parent, "/") + 1 == strings.Count(child, "/") &&
		strings.HasPrefix(child, parent)
}