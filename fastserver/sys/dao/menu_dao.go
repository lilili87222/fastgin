package dao

import (
	"fastgin/config"
	"fastgin/sys/model"
	"github.com/thoas/go-funk"
)

type MenuDao struct {
}

//func NewMenuDao() MenuDao {
//	return MenuDao{}
//}

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
	// parentId为0的是根菜单
	return GenMenuTree(0, menus), err
}

func GenMenuTree(parentId uint, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)

	for _, m := range menus {
		if *m.ParentId == parentId {
			children := GenMenuTree(m.ID, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
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

// 根据用户ID获取用户的权限(可访问)菜单列表
func (m MenuDao) GetUserMenusByUserId(userId uint) ([]*model.Menu, error) {
	// 获取用户
	var user model.User
	err := config.DB.Where("id = ?", userId).Preload("Roles").First(&user).Error
	if err != nil {
		return nil, err
	}
	// 获取角色
	roles := user.Roles
	// 所有角色的菜单集合
	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range roles {
		var userRole model.Role
		err := config.DB.Where("id = ?", role.ID).Preload("Menus").First(&userRole).Error
		if err != nil {
			return nil, err
		}
		// 获取角色的菜单
		menus := userRole.Menus
		allRoleMenus = append(allRoleMenus, menus...)
	}

	// 所有角色的菜单集合去重
	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.ID))
	}
	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.ID) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	// 获取状态status为1的菜单
	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, err
}

// 根据用户ID获取用户的权限(可访问)菜单树
func (m MenuDao) GetUserMenuTreeByUserId(userId uint) ([]*model.Menu, error) {
	menus, err := m.GetUserMenusByUserId(userId)
	if err != nil {
		return nil, err
	}
	tree := GenMenuTree(0, menus)
	return tree, err
}