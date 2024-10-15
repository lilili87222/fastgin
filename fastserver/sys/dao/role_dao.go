package dao

import (
	"fastgin/config"
	"fastgin/sys/dto"
	"fastgin/sys/model"
	"fmt"
	"strings"
)

type RoleDao struct {
}

// 获取角色列表
func (r RoleDao) GetRoles(req *dto.RoleListRequest) ([]model.Role, int64, error) {
	var list []model.Role
	db := config.DB.Model(&model.Role{}).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

// 根据角色ID获取角色
func (r RoleDao) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := config.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// 创建角色
func (r RoleDao) CreateRole(role *model.Role) error {
	err := config.DB.Create(role).Error
	return err
}

// 更新角色
func (r RoleDao) UpdateRoleById(roleId uint, role *model.Role) error {
	err := config.DB.Model(&model.Role{}).Where("id = ?", roleId).Updates(role).Error
	return err
}

// 获取角色的权限菜单
func (r RoleDao) GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	var role model.Role
	err := config.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// 更新角色的权限菜单
func (r RoleDao) UpdateRoleMenus(role *model.Role) error {
	err := config.DB.Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// New function to get roles by IDs
// 删除角色
func (r RoleDao) BatchDeleteRoleByIds(roleIds []uint) ([]*model.Role, error) {
	roles, err := r.GetRolesByIds(roleIds)
	if err != nil {
		return nil, err
	}
	return roles, config.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
}
