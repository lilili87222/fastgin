package sys

import (
	"fastgin/config"
	"fastgin/internal/bean"
	"fastgin/internal/controller"
	sys2 "fastgin/internal/dao/sys"
	"fastgin/internal/model/sys"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"strconv"
)

type IMenuController interface {
	GetMenus(c *gin.Context)             // 获取菜单列表
	GetMenuTree(c *gin.Context)          // 获取菜单树
	CreateMenu(c *gin.Context)           // 创建菜单
	UpdateMenuById(c *gin.Context)       // 更新菜单
	BatchDeleteMenuByIds(c *gin.Context) // 批量删除菜单

	GetUserMenusByUserId(c *gin.Context)    // 获取用户的可访问菜单列表
	GetUserMenuTreeByUserId(c *gin.Context) // 获取用户的可访问菜单树
}

type MenuController struct {
	MenuRepository sys2.IMenuRepository
}

func NewMenuController() IMenuController {
	menuRepository := sys2.NewMenuRepository()
	menuController := MenuController{MenuRepository: menuRepository}
	return menuController
}

// 获取菜单列表
func (mc MenuController) GetMenus(c *gin.Context) {
	menus, err := mc.MenuRepository.GetMenus()
	if err != nil {
		controller.Fail(c, nil, "获取菜单列表失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"menus": menus}, "获取菜单列表成功")
}

// 获取菜单树
func (mc MenuController) GetMenuTree(c *gin.Context) {
	menuTree, err := mc.MenuRepository.GetMenuTree()
	if err != nil {
		controller.Fail(c, nil, "获取菜单树失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"menuTree": menuTree}, "获取菜单树成功")
}

// 创建菜单
func (mc MenuController) CreateMenu(c *gin.Context) {
	var req bean.CreateMenuRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}

	// 获取当前用户
	ur := sys2.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		controller.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	menu := sys.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       &req.Icon,
		Path:       req.Path,
		Redirect:   &req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: &req.ActiveMenu,
		ParentId:   &req.ParentId,
		Creator:    ctxUser.Username,
	}

	err = mc.MenuRepository.CreateMenu(&menu)
	if err != nil {
		controller.Fail(c, nil, "创建菜单失败: "+err.Error())
		return
	}
	controller.Success(c, nil, "创建菜单成功")
}

// 更新菜单
func (mc MenuController) UpdateMenuById(c *gin.Context) {
	type UpdateMenuRequest struct {
		Name       string `json:"name" form:"name" validate:"required,min=1,max=50"`
		Title      string `json:"title" form:"title" validate:"required,min=1,max=50"`
		Icon       string `json:"icon" form:"icon" validate:"min=0,max=50"`
		Path       string `json:"path" form:"path" validate:"required,min=1,max=100"`
		Redirect   string `json:"redirect" form:"redirect" validate:"min=0,max=100"`
		Component  string `json:"component" form:"component" validate:"min=0,max=100"`
		Sort       uint   `json:"sort" form:"sort" validate:"gte=1,lte=999"`
		Status     uint   `json:"status" form:"status" validate:"oneof=1 2"`
		Hidden     uint   `json:"hidden" form:"hidden" validate:"oneof=1 2"`
		NoCache    uint   `json:"noCache" form:"noCache" validate:"oneof=1 2"`
		AlwaysShow uint   `json:"alwaysShow" form:"alwaysShow" validate:"oneof=1 2"`
		Breadcrumb uint   `json:"breadcrumb" form:"breadcrumb" validate:"oneof=1 2"`
		ActiveMenu string `json:"activeMenu" form:"activeMenu" validate:"min=0,max=100"`
		ParentId   uint   `json:"parentId" form:"parentId"`
	}
	var req UpdateMenuRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}

	// 获取路径中的menuId
	menuId, _ := strconv.Atoi(c.Param("menuId"))
	if menuId <= 0 {
		controller.Fail(c, nil, "菜单ID不正确")
		return
	}

	// 获取当前用户
	ur := sys2.NewUserRepository()
	ctxUser, err := ur.GetCurrentUser(c)
	if err != nil {
		controller.Fail(c, nil, "获取当前用户信息失败")
		return
	}

	menu := sys.Menu{
		Name:       req.Name,
		Title:      req.Title,
		Icon:       &req.Icon,
		Path:       req.Path,
		Redirect:   &req.Redirect,
		Component:  req.Component,
		Sort:       req.Sort,
		Status:     req.Status,
		Hidden:     req.Hidden,
		NoCache:    req.NoCache,
		AlwaysShow: req.AlwaysShow,
		Breadcrumb: req.Breadcrumb,
		ActiveMenu: &req.ActiveMenu,
		ParentId:   &req.ParentId,
		Creator:    ctxUser.Username,
	}

	err = mc.MenuRepository.UpdateMenuById(uint(menuId), &menu)
	if err != nil {
		controller.Fail(c, nil, "更新菜单失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "更新菜单成功")

}

// 批量删除菜单
func (mc MenuController) BatchDeleteMenuByIds(c *gin.Context) {
	var req bean.DeleteMenuRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	err := mc.MenuRepository.BatchDeleteMenuByIds(req.MenuIds)
	if err != nil {
		controller.Fail(c, nil, "删除菜单失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "删除菜单成功")
}

// 根据用户ID获取用户的可访问菜单列表
func (mc MenuController) GetUserMenusByUserId(c *gin.Context) {
	// 获取路径中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		controller.Fail(c, nil, "用户ID不正确")
		return
	}

	menus, err := mc.MenuRepository.GetUserMenusByUserId(uint(userId))
	if err != nil {
		controller.Fail(c, nil, "获取用户的可访问菜单列表失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"menus": menus}, "获取用户的可访问菜单列表成功")
}

// 根据用户ID获取用户的可访问菜单树
func (mc MenuController) GetUserMenuTreeByUserId(c *gin.Context) {
	// 获取路径中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		controller.Fail(c, nil, "用户ID不正确")
		return
	}

	menuTree, err := mc.MenuRepository.GetUserMenuTreeByUserId(uint(userId))
	if err != nil {
		controller.Fail(c, nil, "获取用户的可访问菜单树失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"menuTree": menuTree}, "获取用户的可访问菜单树成功")
}
