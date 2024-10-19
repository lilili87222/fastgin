package service

import (
	"fastgin/database"
	"fastgin/modules/sys/dao"
	"fastgin/modules/sys/model"
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
	return database.ListAll[*model.Menu]("sort")
	//return s.menuDao.List()
}

// 获取菜单树
func (s *MenuService) GetMenuTree() ([]*model.Menu, error) {
	menus, err := s.GetMenus()
	if err != nil {
		return nil, err
	}
	return GenMenuTree(0, menus), nil
}

// 创建菜单
func (s *MenuService) CreateMenu(menu *model.Menu) error {
	return database.Create(menu)
}

// 更新菜单
func (s *MenuService) UpdateMenuById(menu *model.Menu) error {
	return database.Update(menu)
	//return s.menuDao.Update(menuId, menu)
}

// 批量删除菜单
func (s *MenuService) BatchDeleteMenuByIds(menuIds []uint) error {
	return s.menuDao.BatchDeleteMenuByIds(menuIds)
}

// 根据用户ID获取用户的权限(可访问)菜单列表
func (s *MenuService) GetUserMenusByUserId(userId uint) ([]*model.Menu, error) {
	userDao := dao.NewUserDao()
	roleDao := dao.RoleDao{}
	user, err := userDao.GetUserWithRoles(userId)
	if err != nil {
		return nil, err
	}

	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range user.Roles {
		//userRole, err := userDao.GetRoleWithMenus(role.Id)
		userRole, err := roleDao.GetRoleWithMenus(role.Id)
		//userRole, err := s.menuDao.GetRoleWithMenus(role.Id)
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

func (s *MenuService) InsertAppMenuToAdmin(menu model.Menu) {
	menuList, _ := database.ListAll[model.Menu]()
	for _, m := range menuList {
		if m.Component == menu.Component {
			return
		}
	}
	roles := []model.Role{
		{
			Model:   model.Model{Id: 1},
			Name:    "管理员",
			Keyword: "admin",
			Desc:    new(string),
			Sort:    1,
			Status:  1,
			Creator: "系统",
		},
	}
	appMenuId := uint(8)
	menu.ParentId = &appMenuId
	menu.Roles = roles
	database.Create(&menu)
}
