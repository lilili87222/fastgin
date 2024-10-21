package controller

import (
	"fastgin/common/httpz"
	"fastgin/config"
	"fastgin/modules/sys/dto"
	"fastgin/modules/sys/model"
	"fastgin/modules/sys/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"strconv"
)

// MenuController handles menu-related requests
type MenuController struct {
	menuService *service.MenuService
	userService *service.UserService
}

// NewMenuController creates a new MenuController
func NewMenuController() *MenuController {
	return &MenuController{menuService: service.NewMenuService(), userService: service.NewUserService()}
}

// List retrieves a list of menus
// @Summary Get menu list
// @Description Get a list of menus
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/menus [get]
func (mc *MenuController) List(c *gin.Context) {
	menus, err := mc.menuService.GetMenus()
	if err != nil {
		httpz.ServerError(c, "获取菜单列表失败: "+err.Error())
		return
	}
	//util.Success(c, gin.H{"menus": menus}, "获取菜单列表成功")
	httpz.Success(c, menus)
}

// GetMenuTree retrieves the menu tree
// @Summary Get menu tree
// @Description Get the menu tree
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/menu/tree [get]
func (mc *MenuController) GetMenuTree(c *gin.Context) {
	menuTree, err := mc.menuService.GetMenuTree()
	if err != nil {
		httpz.ServerError(c, "获取菜单树失败: "+err.Error())
		return
	}
	//util.Success(c, gin.H{"MenuTree": menuTree}, "获取菜单树成功")
	httpz.Success(c, menuTree)
}

// Create creates a new menu
// @Summary Create menu
// @Description Create a new menu
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param menu body dto.CreateMenuRequest true "Create menu request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/menu/index [post]
func (mc *MenuController) Create(c *gin.Context) {
	var req dto.CreateMenuRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		httpz.BadRequest(c, errStr)
		return
	}
	ctxUser, err := mc.userService.GetCurrentUser(c)
	if err != nil {
		httpz.ServerError(c, "获取当前用户信息失败")
		return
	}
	menu := model.Menu{
		//Name:       req.Name,
		//Title:      req.Title,
		//Icon:       &req.Icon,
		//Path:       req.Path,
		//Redirect:   &req.Redirect,
		//Component:  req.Component,
		//Sort:       req.Sort,
		//Status:     req.Status,
		//Hidden:     req.Hidden,
		//NoCache:    req.NoCache,
		//AlwaysShow: req.AlwaysShow,
		//Breadcrumb: req.Breadcrumb,
		//ActiveMenu: &req.ActiveMenu,
		//ParentId:   &req.ParentId,
		Creator: ctxUser.UserName,
	}
	e := copier.Copy(&menu, &req)

	if e != nil {
		httpz.ServerError(c, "创建菜单失败: "+e.Error())
		return
	}
	err = mc.menuService.CreateMenu(&menu)
	if err != nil {
		httpz.ServerError(c, "创建菜单失败: "+err.Error())
		return
	}
	httpz.Success(c, nil)
}

// Update updates an existing menu by ID
// @Summary Update menu
// @Description Update an existing menu by ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param menuId path int true "Menu ID"
// @Param menu body dto.CreateMenuRequest true "Update menu request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/menu/index/{menuId} [put]
func (mc *MenuController) Update(c *gin.Context) {
	var req dto.CreateMenuRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		httpz.BadRequest(c, errStr)
		return
	}
	menuId, _ := strconv.Atoi(c.Param("menuId"))
	if menuId <= 0 {
		httpz.BadRequest(c, "菜单ID不正确")
		return
	}
	//ur := service.NewUserService()
	ctxUser, err := mc.userService.GetCurrentUser(c)
	if err != nil {
		httpz.ServerError(c, "获取当前用户信息失败")
		return
	}
	menu := model.Menu{
		//Name:       req.Name,
		//Title:      req.Title,
		//Icon:       &req.Icon,
		//Path:       req.Path,
		//Redirect:   &req.Redirect,
		//Component:  req.Component,
		//Sort:       req.Sort,
		//Status:     req.Status,
		//Hidden:     req.Hidden,
		//NoCache:    req.NoCache,
		//AlwaysShow: req.AlwaysShow,
		//Breadcrumb: req.Breadcrumb,
		//ActiveMenu: &req.ActiveMenu,
		//ParentId:   &req.ParentId,
		Creator: ctxUser.UserName,
	}
	e := copier.Copy(&menu, &req)
	if e != nil {
		httpz.ServerError(c, "更新菜单失败: "+e.Error())
		return
	}
	menu.ID = uint64(menuId)
	err = mc.menuService.UpdateMenuById(&menu)
	if err != nil {
		httpz.ServerError(c, "更新菜单失败: "+err.Error())
		return
	}
	httpz.Success(c, nil)
}

// BatchDeleteByIds deletes multiple menus by their Ids
// @Summary Batch delete menus
// @Description BatchDelete multiple menus by their Ids
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param menuIds body httpz.IdListRequest true "BatchDelete menu request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/menu/index [delete]
func (mc *MenuController) BatchDeleteByIds(c *gin.Context) {
	var req httpz.IdListRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	err := mc.menuService.BatchDeleteMenuByIds(req.Ids)
	if err != nil {
		httpz.ServerError(c, "删除菜单失败: "+err.Error())
		return
	}
	httpz.Success(c, nil)
}

// GetUserMenusByUserId retrieves the accessible menus for a user by user ID
// @Summary Get user menus by user ID
// @Description Get the accessible menus for a user by user ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param userId path int true "User ID"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/menu/user/{userId} [get]
func (mc *MenuController) GetUserMenusByUserId(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		httpz.BadRequest(c, "用户ID不正确")
		return
	}
	menus, err := mc.menuService.GetUserMenusByUserId(uint64(userId))
	if err != nil {
		httpz.ServerError(c, "获取用户的可访问菜单列表失败: "+err.Error())
		return
	}
	//util.Success(c, gin.H{"Menus": menus}, "获取用户的可访问菜单列表成功")
	httpz.Success(c, menus)
}

// GetUserMenuTreeByUserId retrieves the accessible menu tree for a user by user ID
// @Summary Get user menu tree by user ID
// @Description Get the accessible menu tree for a user by user ID
// @Tags Menu
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param userId path int true "User ID"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/menu/user_tree/{userId} [get]
func (mc *MenuController) GetUserMenuTreeByUserId(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		httpz.BadRequest(c, "用户ID不正确")
		return
	}
	menuTree, err := mc.menuService.GetUserMenuTreeByUserId(uint64(userId))
	if err != nil {
		httpz.ServerError(c, "获取用户的可访问菜单树失败: "+err.Error())
		return
	}
	//util.Success(c, gin.H{"MenuTree": menuTree}, "获取用户的可访问菜单树成功")
	httpz.Success(c, menuTree)
}
