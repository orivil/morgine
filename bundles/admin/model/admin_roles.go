// Copyright 2019 orivil.com. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found at https://mit-license.org.

package admin_model

type AdminRole struct {
	ID int
	AdminID int `gorm:"unique_index:admin_role_idx"`
	RoleID int `gorm:"unique_index:admin_role_idx"`
}

// 统计管理员角色数量
func CountAdminRoles(adminID, roleID int) (total int, err error) {
	model := DB.Model(&AdminRole{})
	if adminID > 0 {
		model = model.Where("admin_id=?", adminID)
	} else if roleID > 0 {
		model = model.Where("role_id=?", roleID)
	}
	err = model.Count(&total).Error
	return
}

type AdminRoleWithName struct {
	ID int
	AdminID int
	RoleID int
	AdminName string
	RoleName string
}

func GetAdminRoleWithNames(adminID, roleID, limit, offset int) (roles []*AdminRoleWithName) {
	model := DB.Model(&AdminRole{})
	if adminID > 0 {
		model = model.Where("admin_id=?", adminID)
	} else if roleID > 0 {
		model = model.Where("role_id=?", roleID)
	}
	var ars []*AdminRole
	model.Order("id").Limit(limit).Offset(offset).Find(&ars)
	if ln := len(ars); ln > 0 {
		var roleIDs= make(map[int]struct{}, ln)
		var adminIDs= make(map[int]struct{}, ln)
		for _, ar := range ars {
			roleIDs[ar.RoleID] = struct{}{}
			adminIDs[ar.AdminID] = struct{}{}
		}
		var rs []*Role
		var as []*Admin
		DB.Where("id in (?)", getSlice(roleIDs)).Find(&rs)
		DB.Where("id in (?)", getSlice(adminIDs)).Find(&as)
		var rsm = make(map[int]*Role, len(rs))
		var asm = make(map[int]*Admin, len(as))
		for _, r := range rs {
			rsm[r.ID] = r
		}
		for _, a := range as {
			asm[a.ID] = a
		}
		roles = make([]*AdminRoleWithName, ln)
		for key, ar := range ars {
			roles[key] = &AdminRoleWithName {
				ID: ar.ID,
				AdminID: ar.AdminID,
				RoleID: ar.RoleID,
				AdminName: asm[ar.AdminID].Username,
				RoleName: rsm[ar.RoleID].Name,
			}
		}
	}
	return
}

// 获得管理员角色列表
func GetAdminRoles(adminID, limit, offset int) (roles []*Role, err error) {
	expr := DB.Model(&AdminRole{}).Where("admin_id=?", adminID).Order("id").Limit(limit).Offset(offset).Select("role_id").QueryExpr()
	err = DB.Where("id in (?)", expr).Find(&roles).Error
	return
}

// 设置管理员角色
func SetAdminRole(adminID, roleID int) error {
	return DB.Create(&AdminRole{AdminID:adminID, RoleID:roleID}).Error
}

// 移除管理员角色
func DelAdminRoles(ids []int) error {
	return DB.Where("id in (?)", ids).Delete(&AdminRole{}).Error
}

// 统计角色管理员数量
func CountRoleAdmins(roleID int) (total int, err error) {
	err = DB.Model(&AdminRole{}).Where("role_id=?", roleID).Count(&total).Error
	return
}

// 获得包含指定角色的管理员列表
func GetRoleAdmins(roleID, limit, offset int) (admins []*Admin, err error) {
	expr := DB.Model(&AdminRole{}).Where("role_id=?", roleID).Order("id").Limit(limit).Offset(offset).Select("admin_id").QueryExpr()
	err = DB.Where("id in (?)", expr).Find(&admins).Error
	return
}

// 获得管理员-角色列表
func GetRoleList(limit, offset int) (ars []*AdminRole, err error) {
	err = DB.Limit(limit).Offset(offset).Find(&ars).Error
	return
}

func getSlice(im map[int]struct{}) []int {
	var s []int
	for key := range im {
		s = append(s, key)
	}
	return s
}