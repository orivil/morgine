// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin_model

type Role struct {
	ID int
	Name string `gorm:"unique_index"`
	AuthIDs string `desc:"权限ID, 如: 1,2"`
	AuthCa string `desc:"角色权限控制器-方法, 如: Goods-list,Goods-add"`
}

func CountRoles() (total int, err error) {
	err = DB.Model(&Role{}).Count(&total).Error
	return
}

func GetRoles(limit, offset int) (roles []*Role, err error) {
	err = DB.Order("id").Limit(limit).Offset(offset).Find(&roles).Error
	return
}

func (r *Role) Create() error {
	r.ID = 0
	return DB.Create(r).Error
}

func (r *Role) Update() error {
	return DB.Model(r).Updates(r).Error
}

func (r *Role) Delete() error {
	return DB.Delete(r).Error
}

func DeleteRoles(ids []int) error {
	return DB.Where("id in (?)", ids).Delete(&Role{}).Error
}