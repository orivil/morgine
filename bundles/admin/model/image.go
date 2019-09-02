// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin_model

import (
	"errors"
	"github.com/orivil/morgine/bundles/utils/sql"
	"os"
	"time"
)

var ErrLabelIsNotEmpty = errors.New("image label contains images")
var ErrLabelKeyIsExist = errors.New("image label key already exist")

type Image struct {
	ID int
	AdminID int `gorm:"index"`
	File string `gorm:"index"`
	Url string `gorm:"-" desc:"服务地址"`
	LabelKey string `gorm:"index"`
	CreatedAt *time.Time `gorm:"index"`
}

func CountLabelImage(adminID int, labelKey string) (totalImage int) {
	DB.Model(&Image{}).Where("admin_id=? AND label_key=?", adminID, labelKey).Count(&totalImage)
	return
}

func GetImageByFile(file string) *Image {
	exist := &Image{}
	DB.Where("file=?", file).First(exist)
	if exist.ID > 0 {
		return exist
	} else {
		return nil
	}
}

func GetImages(labelKey string, adminID, limit, offset int) (imgs []*Image) {
	DB.Where("admin_id=? AND label_key=?", adminID, labelKey).Limit(limit).Offset(offset).Order("id desc").Find(&imgs)
	return imgs
}

func (i *Image) Create(adminID int) error {
	i.AdminID = adminID
	return DB.Create(i).Error
}

func DeleteImagesByFiles(files []string) error {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			if !os.IsNotExist(err) {
				return err
			}
		}
	}
	return DB.Where("file in (?)", files).Delete(&Image{}).Error
}

func DeleteImages(adminID int, ids []int) error {
	//var exists []*Image
	//DB.Where("admin_id=? AND id in (?)", adminID, ids).Find(&exists)
	//for _, exist := range exists {
	//	err := os.Remove(exist.File)
	//	if err != nil {
	//		// 如果文件不存在则继续删除
	//		if !os.IsNotExist(err) {
	//			return err
	//		}
	//	}
	//}
	return DB.Where("admin_id=? AND id in (?)", adminID, ids).Delete(&Image{}).Error
}

func GetImagesByIDs(adminID int, ids []int) []*Image {
	var exists []*Image
	DB.Where("admin_id=? AND id in (?)", adminID, ids).Find(&exists)
	return exists
}

type Label struct {
	ID int
	Name string
	Key string `gorm:"unique_index"`
	Width int
	Height int
	AllowedEdit sql.Boolean `desc:"允许后台操作图片(增删改), 1-允许 2-不允许"`
	SizeKB int `desc:"图片大小/KB"`
}

func GetLabel(id int) *Label {
	l := &Label{}
	DB.Model(&Label{ID: id}).Where("id=?", id).First(l)
	if l.ID > 0 {
		return l
	} else {
		return nil
	}
}

func GetLabels() (ls []*Label) {
	DB.Find(&ls)
	return ls
}

func (l *Label) Create(key string) error {
	if l.AllowedEdit == 0 {
		l.AllowedEdit = sql.False
	}
	exist := GetLabelByKey(key)
	if exist != nil {
		return ErrLabelKeyIsExist
	} else {
		l.Key = key
	}
	return DB.Create(l).Error
}

func (l *Label) Update() error {
	return DB.Model(l).Where("id=?", l.ID).Updates(l).Error
}

func (l *Label) Delete(adminID int) error {
	totalImage := CountLabelImage(adminID, l.Key)
	if totalImage > 0 {
		return ErrLabelIsNotEmpty
	} else {
		return DB.Delete(l).Error
	}
}

func GetLabelByKey(key string) *Label {
	exist := &Label{}
	DB.Where("key=?", key).First(exist)
	if exist.ID > 0 {
		return exist
	} else {
		return nil
	}
}

func InitLabel(labels ... *Label) {
	for _, label := range labels {
		exist := GetLabelByKey(label.Key)
		if exist == nil {
			err := label.Create(label.Key)
			if err != nil {
				panic(err)
			}
		} else {
			*label = *exist
		}
	}
}