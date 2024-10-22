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
func (s *MenuService) BatchDeleteMenuByIds(menuIds []uint64) error {
	return s.menuDao.BatchDeleteMenuByIds(menuIds)
}

// 根据用户ID获取用户的权限(可访问)菜单列表
func (s *MenuService) GetUserMenusByUserId(userId uint64) ([]*model.Menu, error) {
	userDao := dao.NewUserDao()
	roleDao := dao.RoleDao{}
	user, err := userDao.GetUserWithRoles(userId)
	if err != nil {
		return nil, err
	}

	allRoleMenus := make([]*model.Menu, 0)
	for _, role := range user.Roles {
		//userRole, err := userDao.GetRoleWithMenus(role.ID)
		userRole, err := roleDao.GetRoleWithMenus(role.ID)
		//userRole, err := s.menuDao.GetRoleWithMenus(role.ID)
		if err != nil {
			return nil, err
		}
		allRoleMenus = append(allRoleMenus, userRole.Menus...)
	}

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

	accessMenus := make([]*model.Menu, 0)
	for _, menu := range allRoleMenusUniq {
		if menu.Status == 1 {
			accessMenus = append(accessMenus, menu)
		}
	}

	return accessMenus, nil
}

// 根据用户ID获取用户的权限(可访问)菜单树
func (s *MenuService) GetUserMenuTreeByUserId(userId uint64) ([]*model.Menu, error) {
	menus, err := s.GetUserMenusByUserId(userId)
	if err != nil {
		return nil, err
	}
	return GenMenuTree(0, menus), nil
}

func GenMenuTree(parentId uint64, menus []*model.Menu) []*model.Menu {
	tree := make([]*model.Menu, 0)
	for _, m := range menus {
		if m.ParentID == parentId {
			children := GenMenuTree(m.ID, menus)
			m.Children = children
			tree = append(tree, m)
		}
	}
	return tree
}

// InsertAppMenuToAdmin 插入应用菜单到管理员, 用于初始化应用菜单, 仅在初始化时调用, 之后不再调用, 如果一键存在，则不会重复插入
func (s *MenuService) InsertAppMenuToAdmin(menu model.Menu) {
	menuList, _ := database.ListAll[model.Menu]()
	for _, m := range menuList {
		if m.Component == menu.Component {
			return
		}
	}
	roles := []model.Role{
		{
			ID:      1,
			Name:    "管理员",
			Keyword: "admin",
			Des:     "",
			Sort:    1,
			Status:  1,
			Creator: "系统",
		},
	}
	//appMenuId := uint64(8)
	//menu.ParentID = 8
	menu.Roles = roles
	database.Create(&menu)
}
