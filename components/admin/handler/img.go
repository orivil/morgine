// Copyright 2020 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package handler

import (
	"errors"
	"github.com/orivil/morgine/components/admin/auth"
	"github.com/orivil/morgine/components/admin/env"
	"github.com/orivil/morgine/components/admin/utils"
	"github.com/orivil/morgine/param"
	"github.com/orivil/morgine/utils/random"
	"github.com/orivil/morgine/xx"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func GetImageDirs(method, route string, cdt *xx.Condition) {
	doc := &xx.Doc {
		Title:     "获得图片目录列表",
		Desc:      "",
		Params:    nil,
		Responses: xx.Responses{
			{
				Body: xx.JsonData(xx.StatusSuccess, xx.MAP{
					env.Config.ImgDir: xx.MAP{
						"avatar": xx.MAP{
							"admin": xx.MAP{},
							"user": xx.MAP{},
						},
						"photo": xx.MAP{},
					},
				}),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		root, err := adminImgDir(ctx)
		if err != nil {
			xx.HandleError(ctx, err)
		} else {
			dir, err := utils.WalkDirs(root)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				dir.Trim(root)
				xx.SendJson(ctx, xx.StatusSuccess, dir)
			}
		}
	})
}

func adminImgDir(ctx *xx.Context) (string, error) {
	id := auth.GetAdminID(ctx)
	if id > 0 {
		return utils.CleanDir(env.Config.ImgDir) + "/admins/" + strconv.Itoa(id), nil
	} else {
		return "", errors.New("用户未登录")
	}
}

type img struct {
	Name string
	Size int64
}

func GetDirImages(method, route string, cdt *xx.Condition) {
	type query struct {
		Dir string `param:"dir" desc:"示例：images/avatar/admin"`
	}
	doc := &xx.Doc{
		Title:     "根据目录获得图片列表",
		Desc:      "同一目录不可上传太多图片",
		Params:    xx.Params{
			{
				Type:xx.Query,
				Schema:&query{},
			},
		},
		Responses: xx.Responses{
			{
				Description: "Size 单位: KB",
				Body: xx.JsonData(xx.StatusSuccess, []*img{{Name: "1.jpg", Size: 60}, {Name: "2.jpg", Size: 1000}}),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		root, err := adminImgDir(ctx)
		if err != nil {
			xx.HandleError(ctx, err)
		} else {
			ps := &query{}
			err := ctx.Unmarshal(ps)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				dir := root + "/" + utils.CleanDir(ps.Dir)
				fs, err := ioutil.ReadDir(dir)
				if err != nil {
					xx.HandleError(ctx, err)
				} else {
					var files []*img
					for _, f := range fs {
						if !f.IsDir() {
							files = append(files, &img{
								Name: dir + "/" + f.Name(),
								Size: f.Size() >> 10,
							})
						}
					}
					xx.SendJson(ctx, xx.StatusSuccess, files)
				}
			}
		}
	})
}

func CreateDir(method, route string, cdt *xx.Condition) {
	type params struct {
		Dir string `param:"dir" required:"" desc:"目录名称, 如 avatar/admin"`
	}
	doc := &xx.Doc{
		Title:     "创建图片目录",
		Params:    xx.Params{
			{
				Type:xx.Form,
				Schema:&params{},
			},
		},
		Responses: xx.Responses{
			{
				Body: xx.Message{},
			},
			{
				Description: "返回父目录下的所有子目录",
				Body: xx.JsonData(xx.StatusSuccess, &utils.Dir{}),
			},
		},
	}

	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		root, err := adminImgDir(ctx)
		if err != nil {
			xx.HandleError(ctx, err)
		} else {
			ps := &params{}
			err = ctx.Unmarshal(ps)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				err = os.Mkdir(filepath.Join(root, ps.Dir), os.ModePerm)
				if err != nil {
					xx.HandleError(ctx, err)
				} else {
					path := utils.CleanDir(ps.Dir)
					d := &utils.Dir {
						Path: path,
						Name: filepath.Base(path),
						Subs: nil,
					}
					xx.SendJson(ctx, xx.StatusSuccess, d)
				}
			}
		}
	})
}

func DelDir(method, route string, cdt *xx.Condition) {
	type params struct {
		Dir string `param:"dir" desc:"需要删除的目录，如 images/avatar/admin"`
	}
	doc := &xx.Doc{
		Title:     "删除目录",
		Desc:      "如果该目录或其子目录下仍有图片则不可删除，不可删除图片根目录",
		Params:    xx.Params{
			{Type:xx.Query, Schema: &params{}},
		},
		Responses: xx.Responses{
			{Body: xx.Message{}},
			{Body: xx.JsonData(xx.StatusSuccess, nil)},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		root, err := adminImgDir(ctx)
		if err != nil {
			xx.HandleError(ctx, err)
		} else {
			ps := &params{}
			err = ctx.Unmarshal(ps)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				// 检查目录及其子目录是否为空目录
				existFile := false
				dir := filepath.Join(root, ps.Dir)
				err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
					if existFile {
						return filepath.SkipDir
					}
					if err != nil {
						return err
					} else {
						if !info.IsDir() {
							existFile = true
							return filepath.SkipDir
						}
					}
					return nil
				})
				if err != nil {
					xx.HandleError(ctx, err)
				} else {
					if existFile {
						xx.SendMessage(ctx, xx.MsgTypeError, "该目录或其子目录下仍有文件")
					} else {
						err = os.Remove(dir)
						if err != nil {
							xx.HandleError(ctx, err)
						} else {
							xx.SendJson(ctx, xx.StatusSuccess, nil)
						}
					}
				}
			}
		}
	})
}

func UploadDirImage(method, route string, cdt *xx.Condition) {
	type query struct {
		Dir string `param:"dir" desc:"上传的目录，如 images/avatar"`
	}
	type form struct {
		Image param.FileHandler `param:"image" desc:"图片字段"`
	}
	doc := &xx.Doc{
		Title:     "上传图片到目录",
		Desc:      "dir 为 query 参数，image 为 form 参数",
		Params:    xx.Params{
			{
				Type:xx.Query,
				Schema:&query{},
			},
			{
				Type:xx.Form,
				Schema:&form{},
			},
		},
		Responses: xx.Responses{
			{
				Description: "上传成功",
				Body:xx.JsonData(xx.StatusSuccess, &img{
					Name: "",
					Size: 0,
				}),
			},
		},
	}
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		root, err := adminImgDir(ctx)
		if err != nil {
			xx.HandleError(ctx, err)
		} else {
			q := &query{}
			var image *img
			f := &form {
				Image: func(field string, header *multipart.FileHeader) error {
					var file multipart.File
					file, err = header.Open()
					if err != nil {
						return err
					}
					var data []byte
					data, err = ioutil.ReadAll(file)
					if err != nil {
						return err
					}
					ext := filepath.Ext(header.Filename)
					name := string(random.NewRandByte(32))
					filename := filepath.Join(root, q.Dir, name + ext)
					image = &img {
						Name: filename,
						Size: header.Size >> 10,
					}
					return ioutil.WriteFile(filename, data, os.ModePerm)
				},
			}
			err := ctx.Unmarshal(q, f)
			if err != nil {
				xx.HandleError(ctx, err)
			} else {
				xx.SendJson(ctx, xx.StatusSuccess, image)
			}
		}
	})
}

func ServeImage(method string, cdt *xx.Condition) {
	doc := &xx.Doc {
		Title:     "图片服务",
		Desc:      "",
		Params:    nil,
		Responses: nil,
	}
	route := "/" + filepath.Clean(env.Config.ImgDir) + "/"
	cdt.Handle(method, route, doc, func(ctx *xx.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, strings.TrimPrefix(ctx.Request.URL.Path, "/"))
	})
}