package dao

import (
	"fastgin/config"
	"fastgin/sys/model"
)

type MenuDao struct {
}

// 获取菜单列表
func (m MenuDao) GetMenus() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := config.DB.Order("sort").Find(&menus).Error
	return menus, err
}

// 获取菜单树
func (m MenuDao) GetMenuTree() ([]*model.Menu, error) {
	var menus []*model.Menu
	err := config.DB.Order("sort").Find(&menus).Error
	return menus, err
}

// 创建菜单
func (m MenuDao) CreateMenu(menu *model.Menu) error {
	err := config.DB.Create(menu).Error
	return err
}

// 更新菜单
func (m MenuDao) UpdateMenuById(menuId uint, menu *model.Menu) error {
	err := config.DB.Model(menu).Where("id = ?", menuId).Updates(menu).Error
	return err
}

// 批量删除菜单
func (m MenuDao) BatchDeleteMenuByIds(menuIds []uint) error {
	var menus []*model.Menu
	err := config.DB.Where("id IN (?)", menuIds).Find(&menus).Error
	if err != nil {
		return err
	}
	err = config.DB.Select("Roles").Unscoped().Delete(&menus).Error
	return err
}

// 根据用户ID获取用户
func (m MenuDao) GetUserById(userId uint) (*model.User, error) {
	var user model.User
	err := config.DB.Where("id = ?", userId).Preload("Roles").First(&user).Error
	return &user, err
}

// 根据角色ID获取角色
func (m MenuDao) GetRoleById(roleId uint) (*model.Role, error) {
	var role model.Role
	err := config.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return &role, err
}
