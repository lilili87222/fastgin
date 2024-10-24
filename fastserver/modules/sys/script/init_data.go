package script

import (
	"errors"
	config2 "fastgin/boost/config"
	"fastgin/common/util"
	"fastgin/database"
	"fastgin/modules/sys/model"
	"gorm.io/gorm"
	"slices"
)

// 初始化mysql数据
func InitSysModuleDatabase() {
	tableNames, e := database.GetTableNames(database.DB, config2.Configs.Database.MysqlConfig.Database)
	if e != nil {
		config2.Log.Errorf("获取数据库表名失败：%v", e)
		panic(e)
	}
	insertData := false
	tableList := []database.ITableModel{&model.User{}, &model.Role{}, &model.Menu{}, &model.Api{}, &model.OperationLog{}, &model.Dictionary{}}
	for _, tableModel := range tableList {
		if slices.Contains(tableNames, tableModel.TableName()) {
			continue
		} else {
			if e := database.DB.AutoMigrate(tableModel); e != nil {
				config2.Log.Errorf("初始化数据库表失败：%v", e)
				panic(e)
			} else {
				insertData = true
			}
		}
	}
	config2.Log.Infof("初始化数据库完成!")

	if !insertData {
		return
	}
	// 是否初始化数据
	// 1.写入角色数据
	newRoles := make([]model.Role, 0)
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
		{
			ID:      2,
			Name:    "普通用户",
			Keyword: "user",
			Des:     "",
			Sort:    3,
			Status:  1,
			Creator: "系统",
		},
		{
			ID:      3,
			Name:    "访客",
			Keyword: "guest",
			Des:     "",
			Sort:    5,
			Status:  1,
			Creator: "系统",
		},
	}

	for _, role := range roles {
		err := database.DB.First(role, role.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newRoles = append(newRoles, role)
		}
	}

	if len(newRoles) > 0 {
		err := database.DB.Create(newRoles).Error
		if err != nil {
			config2.Log.Errorf("写入系统角色数据失败：%v", err)
		}
	}

	// 2写入菜单
	newMenus := make([]model.Menu, 0)
	//var uint0 uint = 0
	//var uint1 uint = 1
	componentStr := "component"
	systemUserStr := "/system/user"
	userStr := "user"
	peoplesStr := "peoples"
	treeTableStr := "tree-table"
	treeStr := "tree"
	exampleStr := "example"
	logOperationStr := "/log/operation-log"
	//appOperationStr := "/app/operation"
	documentationStr := "documentation"
	//var uint6 uint = 6
	menus := []model.Menu{
		{
			ID:        1,
			Name:      "System",
			Title:     "系统管理",
			Icon:      componentStr,
			Path:      "/system",
			Component: "Layout",
			Redirect:  systemUserStr,
			Sort:      10,
			ParentID:  0,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			ID:        2,
			Name:      "User",
			Title:     "用户管理",
			Icon:      userStr,
			Path:      "user",
			Component: "/system/user/index",
			Sort:      11,
			ParentID:  1,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			ID:        3,
			Name:      "Role",
			Title:     "角色管理",
			Icon:      peoplesStr,
			Path:      "role",
			Component: "/system/role/index",
			Sort:      12,
			ParentID:  1,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			ID:        4,
			Name:      "Menu",
			Title:     "菜单管理",
			Icon:      treeTableStr,
			Path:      "menu",
			Component: "/system/menu/index",
			Sort:      13,
			ParentID:  1,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			ID:        5,
			Name:      "Api",
			Title:     "接口管理",
			Icon:      treeStr,
			Path:      "api",
			Component: "/system/api/index",
			Sort:      14,
			ParentID:  1,
			Roles:     roles[:1],
			Creator:   "系统",
		},
		{
			ID:        6,
			Name:      "Log",
			Title:     "日志管理",
			Icon:      exampleStr,
			Path:      "/log",
			Component: "Layout",
			Redirect:  logOperationStr,
			Sort:      20,
			ParentID:  0,
			Roles:     roles[:2],
			Creator:   "系统",
		},
		{
			ID:        7,
			Name:      "OperationLog",
			Title:     "操作日志",
			Icon:      documentationStr,
			Path:      "operation-log",
			Component: "/log/operation-log/index",
			Sort:      21,
			ParentID:  6,
			Roles:     roles[:2],
			Creator:   "系统",
		}, {
			ID:        8,
			Name:      "App",
			Title:     "应用功能",
			Icon:      "star",
			Path:      "/app",
			Component: "Layout",
			Redirect:  "/app/operation",
			Sort:      20,
			ParentID:  0,
			Roles:     roles[:2],
			Creator:   "系统",
		},
	}
	for _, menu := range menus {
		err := database.DB.First(menu, menu.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newMenus = append(newMenus, menu)
		}
	}
	if len(newMenus) > 0 {
		err := database.DB.Create(newMenus).Error
		if err != nil {
			config2.Log.Errorf("写入系统菜单数据失败：%v", err)
		}
	}

	// 3.写入用户
	newUsers := make([]model.User, 0)
	users := []model.User{
		{
			ID:       1,
			UserName: "admin@admin.com",
			Password: util.GenPasswd("123456"),
			Mobile:   "18888888888",
			Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
			NickName: "",
			Des:      "",
			Status:   1,
			Creator:  "系统",
			Roles:    roles[:1],
		},
		//{
		//	ID:       2,
		//	UserName: "faker",
		//	Password: util.GenPasswd("123456"),
		//	Mobile:   "19999999999",
		//	Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		//	NickName: "",
		//	Des:      "",
		//	Status:   1,
		//	Creator:  "系统",
		//	Roles:    roles[:2],
		//},
		//{
		//	ID:       3,
		//	UserName: "nike",
		//	Password: util.GenPasswd("123456"),
		//	Mobile:   "13333333333",
		//	Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		//	NickName: "",
		//	Des:      "",
		//	Status:   1,
		//	Creator:  "系统",
		//	Roles:    roles[1:2],
		//},
		//{
		//	ID:       4,
		//	UserName: "bob",
		//	Password: util.GenPasswd("123456"),
		//	Mobile:   "15555555555",
		//	Avatar:   "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
		//	NickName: "",
		//	Des:      "",
		//	Status:   1,
		//	Creator:  "系统",
		//	Roles:    roles[2:3],
		//},
	}

	for _, user := range users {
		err := database.DB.First(user, user.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUsers = append(newUsers, user)
		}
	}

	if len(newUsers) > 0 {
		err := database.DB.Create(newUsers).Error
		if err != nil {
			config2.Log.Errorf("写入用户数据失败：%v", err)
		}
	}

	// 4.写入api
	apis := []model.Api{
		//{
		//	Method:   "POST",
		//	Path:     "/base/login",
		//	Category: "base",
		//	Des:     "用户登录",
		//	Creator:  "系统",
		//},
		{
			Method:   "POST",
			Path:     "/user/logout",
			Category: "user",
			Des:      "用户登出",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/user/refreshToken",
			Category: "user",
			Des:      "刷新JWT令牌",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/user/info",
			Category: "user",
			Des:      "获取当前登录用户信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/user/index",
			Category: "user",
			Des:      "获取用户列表",
			Creator:  "系统",
		},
		{
			Method:   "PUT",
			Path:     "/user/changePwd",
			Category: "user",
			Des:      "更新用户登录密码",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/user/index",
			Category: "user",
			Des:      "创建用户",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/user/index/:userId",
			Category: "user",
			Des:      "更新用户",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/user/index",
			Category: "user",
			Des:      "批量删除用户",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/index",
			Category: "role",
			Des:      "获取角色列表",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/role/index",
			Category: "role",
			Des:      "创建角色",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/role/index/:roleId",
			Category: "role",
			Des:      "更新角色",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/menus/:roleId",
			Category: "role",
			Des:      "获取角色的权限菜单",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/role/menus/:roleId",
			Category: "role",
			Des:      "更新角色的权限菜单",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/role/apis/:roleId",
			Category: "role",
			Des:      "获取角色的权限接口",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/role/apis/:roleId",
			Category: "role",
			Des:      "更新角色的权限接口",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/role/index",
			Category: "role",
			Des:      "批量删除角色",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/index",
			Category: "menu",
			Des:      "获取菜单列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/tree",
			Category: "menu",
			Des:      "获取菜单树",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/menu/index",
			Category: "menu",
			Des:      "创建菜单",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/menu/index/:menuId",
			Category: "menu",
			Des:      "更新菜单",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/menu/index",
			Category: "menu",
			Des:      "批量删除菜单",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/user/:userId",
			Category: "menu",
			Des:      "获取用户的可访问菜单列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/menu/user_tree/:userId",
			Category: "menu",
			Des:      "获取用户的可访问菜单树",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/api/index",
			Category: "api",
			Des:      "获取接口列表",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/api/tree",
			Category: "api",
			Des:      "获取接口树",
			Creator:  "系统",
		},
		{
			Method:   "POST",
			Path:     "/api/index",
			Category: "api",
			Des:      "创建接口",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/api/index/:roleId",
			Category: "api",
			Des:      "更新接口",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/api/index",
			Category: "api",
			Des:      "批量删除接口",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/log/index",
			Category: "log",
			Des:      "获取操作日志列表",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/log/index",
			Category: "log",
			Des:      "批量删除操作日志",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/system/info",
			Category: "system",
			Des:      "获取系统信息",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/system/stop",
			Category: "system",
			Des:      "停止web服务",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/system/restart",
			Category: "system",
			Des:      "重启web服务",
			Creator:  "系统",
		},
	}
	newApi := make([]model.Api, 0)
	newRoleCasbin := make([]model.RoleCasbin, 0)
	for i, api := range apis {
		api.ID = uint64(i + 1)
		err := database.DB.First(api, api.ID).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newApi = append(newApi, api)

			// 管理员拥有所有API权限
			newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
				Keyword: roles[0].Keyword,
				Path:    api.Path,
				Method:  api.Method,
			})

			// 非管理员拥有基础权限
			basePaths := []string{
				//"/base/login",
				"/user/logout",
				"/user/refreshToken",
				"/user/info",
				"/menu/user_tree/:userId",
			}

			if slices.Contains(basePaths, api.Path) {
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[1].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
				newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
					Keyword: roles[2].Keyword,
					Path:    api.Path,
					Method:  api.Method,
				})
			}
		}
	}

	if len(newApi) > 0 {
		if err := database.DB.Create(newApi).Error; err != nil {
			config2.Log.Errorf("写入api数据失败：%v", err)
		}
	}

	if len(newRoleCasbin) > 0 {
		rules := make([][]string, 0)
		for _, c := range newRoleCasbin {
			rules = append(rules, []string{
				c.Keyword, c.Path, c.Method,
			})
		}
		isAdd, err := config2.CasbinEnforcer.AddPolicies(rules)
		if !isAdd {
			config2.Log.Errorf("写入casbin数据失败：%v", err)
		}
	}
}
