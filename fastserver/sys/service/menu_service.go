package service

import (
	"fastgin/sys/dao"
	"fastgin/sys/model"
)

type MenuService struct {
	menuDao dao.MenuDao
}

func NewMenuService() MenuService {
	return MenuService{
		menuDao: dao.MenuDao{},
	}
}

// 获取菜单列表
func (s MenuService) GetMenus() ([]*model.Menu, error) {
	return s.menuDao.GetMenus()
}

// 获取菜单树
func (s MenuService) GetMenuTree() ([]*model.Menu, error) {
	return s.menuDao.GetMenuTree()
}

// 创建菜单
func (s MenuService) CreateMenu(menu *model.Menu) error {
	return s.menuDao.CreateMenu(menu)
}

// 更新菜单
func (s MenuService) UpdateMenuById(menuId uint, menu *model.Menu) error {
	return s.menuDao.UpdateMenuById(menuId, menu)
}

// 批量删除菜单
func (s MenuService) BatchDeleteMenuByIds(menuIds []uint) error {
	return s.menuDao.BatchDeleteMenuByIds(menuIds)
}

// 根据用户ID获取用户的权限(可访问)菜单列表
func (s MenuService) GetUserMenusByUserId(userId uint) ([]*model.Menu, error) {
	return s.menuDao.GetUserMenusByUserId(userId)
}

// 根据用户ID获取用户的权限(可访问)菜单树
func (s MenuService) GetUserMenuTreeByUserId(userId uint) ([]*model.Menu, error) {
	return s.menuDao.GetUserMenuTreeByUserId(userId)
}
