package dao

import (
	"fastgin/database"
	"fastgin/sys/model"
)

type RoleDao struct {
}

// 获取角色的权限菜单
func (r *RoleDao) GetRoleWithMenus(roleId uint) (*model.Role, error) {
	return database.GetByIdPreload[*model.Role](roleId, "Menus")
}
func (r *RoleDao) GetRoleWithUsers(roleId uint) (model.Role, error) {
	return database.GetByIdPreload[model.Role](roleId, "Users")
}

// 更新角色的权限菜单
func (r *RoleDao) UpdateRoleMenus(role *model.Role) error {
	return database.DB.Model(role).Association("Menus").Replace(role.Menus)
}

// 删除角色
func (r *RoleDao) BatchDeleteRoleByIds(roleIds []uint) ([]model.Role, error) {
	roles, err := database.GetByIds[model.Role](roleIds)
	if err != nil {
		return nil, err
	}
	return roles, database.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
}
