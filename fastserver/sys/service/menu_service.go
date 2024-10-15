package service

import (
	"fastgin/sys/dao"
	"fastgin/sys/model"
	"github.com/thoas/go-funk"
)

type MenuService struct {
	menuDao *dao.MenuDao
}

func NewMenuService() *MenuService {
	return &MenuService{
		menuDao: &dao.MenuDao{},
	}
}

// 获取菜单列表
func (s *MenuService) GetMenus() ([]*model.Menu, error) {
	return s.menuDao.GetMenus()
}

// 获取菜单树
func (s *MenuService) GetMenuTree() ([]*model.Menu, error) {
	menus, err := s.menuDao.GetMenuTree()
	if err != nil {
		return nil, err
	}
	return GenMenuTree(0, menus), nil
}

// 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	return s.menuDao.CreateMenu(menu)
}

// 更新菜单
func (s *MenuService) UpdateMenuById(menuId uint, menu *model.Menu) error {
	return s.menuDao.UpdateMenuById(menuId, menu)
}

// 批量删除菜单
func (s *MenuService) BatchDeleteMenuByIds(menuIds []uint) error {
	return s.menuDao.BatchDeleteMenuByIds(menuIds)
}

// 根据用户ID获取用户的权限(可访问)菜单列表
func (s *MenuService) GetUserMenusByUserId(userId uint) ([]*model.Menu, error) {
	user, err := s.menuDao.GetUserById(userId)
	if err != nil {
		return nil, err
	}

	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range user.Roles {
		userRole, err := s.menuDao.GetRoleById(role.Id)
		if err != nil {
			return nil, err
		}
		allRoleMenus = append(allRoleMenus, userRole.Menus...)
	}

	allRoleMenusId := make([]int, 0)
	for _, menu := range allRoleMenus {
		allRoleMenusId = append(allRoleMenusId, int(menu.Id))
	}
	allRoleMenusIdUniq := funk.UniqInt(allRoleMenusId)
	allRoleMenusUniq := make([]*model.Menu, 0)
	for _, id := range allRoleMenusIdUniq {
		for _, menu := range allRoleMenus {
			if id == int(menu.Id) {
				allRoleMenusUniq = append(allRoleMenusUniq, menu)
				break
			}
		}
	}

	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, nil
}

// 根据用户ID获取用户的权限(可访问)菜单树
func (s *MenuService) GetUserMenuTreeByUserId(userId uint) ([]*model.Menu, error) {
	menus, err := s.GetUserMenusByUserId(userId)
	if err != nil {
		return nil, err
	}
	return GenMenuTree(0, menus), nil
}

func GenMenuTree(parentId uint, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)
	for _, m := range menus {
		if *m.ParentId == parentId {
			children := GenMenuTree(m.Id, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}
