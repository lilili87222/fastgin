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

// MenuController handles menu-related requests
type MenuController struct {
	MenuRepository sys2.MenuRepository
}

// NewMenuController creates a new MenuController
func NewMenuController() MenuController {
	menuRepository := sys2.NewMenuRepository()
	menuController := MenuController{MenuRepository: menuRepository}
	return menuController
}

// GetMenus retrieves a list of menus
// @Summary Get menu list
// @Description Get a list of menus
// @Tags Menu
// @Accept json
// @Produce json
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/menus [get]
func (mc MenuController) GetMenus(c *gin.Context) {
	menus, err := mc.MenuRepository.GetMenus()
	if err != nil {
		controller.Fail(c, nil, "获取菜单列表失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"menus": menus}, "获取菜单列表成功")
}

// GetMenuTree retrieves the menu tree
// @Summary Get menu tree
// @Description Get the menu tree
// @Tags Menu
// @Accept json
// @Produce json
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/menu/tree [get]
func (mc MenuController) GetMenuTree(c *gin.Context) {
	menuTree, err := mc.MenuRepository.GetMenuTree()
	if err != nil {
		controller.Fail(c, nil, "获取菜单树失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"menuTree": menuTree}, "获取菜单树成功")
}

// CreateMenu creates a new menu
// @Summary Create menu
// @Description Create a new menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param menu body bean.CreateMenuRequest true "Create menu request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/menu [post]
func (mc MenuController) CreateMenu(c *gin.Context) {
	var req bean.CreateMenuRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
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

// UpdateMenuById updates an existing menu by ID
// @Summary Update menu
// @Description Update an existing menu by ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param menuId path int true "Menu ID"
// @Param menu body bean.UpdateMenuRequest true "Update menu request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/menu/{menuId} [put]
func (mc MenuController) UpdateMenuById(c *gin.Context) {
	var req bean.UpdateMenuRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	menuId, _ := strconv.Atoi(c.Param("menuId"))
	if menuId <= 0 {
		controller.Fail(c, nil, "菜单ID不正确")
		return
	}
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

// BatchDeleteMenuByIds deletes multiple menus by their IDs
// @Summary Batch delete menus
// @Description Delete multiple menus by their IDs
// @Tags Menu
// @Accept json
// @Produce json
// @Param menuIds body bean.DeleteMenuRequest true "Delete menu request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/menu/batch_delete [delete]
func (mc MenuController) BatchDeleteMenuByIds(c *gin.Context) {
	var req bean.DeleteMenuRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
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

// GetUserMenusByUserId retrieves the accessible menus for a user by user ID
// @Summary Get user menus by user ID
// @Description Get the accessible menus for a user by user ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/{userId}/menus [get]
func (mc MenuController) GetUserMenusByUserId(c *gin.Context) {
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

// GetUserMenuTreeByUserId retrieves the accessible menu tree for a user by user ID
// @Summary Get user menu tree by user ID
// @Description Get the accessible menu tree for a user by user ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/{userId}/menu_tree [get]
func (mc MenuController) GetUserMenuTreeByUserId(c *gin.Context) {
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