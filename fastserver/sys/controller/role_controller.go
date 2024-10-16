package controller

import (
	"fastgin/config"
	"fastgin/database"
	"fastgin/sys/dto"
	"fastgin/sys/model"
	"fastgin/sys/service"
	"fastgin/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
	"slices"
	"strconv"
)

// RoleController handles role-related requests
type RoleController struct {
	roleService *service.RoleService
	userService *service.UserService
}

// NewRoleController creates a new RoleController
func NewRoleController() *RoleController {
	return &RoleController{roleService: service.NewRoleService(), userService: service.NewUserService()}
}

// GetRoles retrieves a list of roles
// @Summary Get role list
// @Description Get a list of roles
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param name query string false "Role name"
// @Param status query int false "Role status"
// @Param pageNum query int false "Page number"
// @Param pageSize query int false "Page size"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/list [get]
func (rc *RoleController) GetRoles(c *gin.Context) {
	params, e := util.GetFormData(c)
	if e != nil {
		util.Fail(c, nil, e.Error())
		return
	}
	data, total, err := database.SearchTable[model.Role](dto.NewSearchRequest(params))
	if err != nil {
		util.Fail(c, nil, "获取角色列表失败: "+err.Error())
		return
	}
	util.Success(c, gin.H{"Data": data, "Total": total}, "获取角色列表成功")
}

// CreateRole creates a new role
// @Summary Create role
// @Description Create a new role
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param role body dto.CreateRoleRequest true "Create role request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/role [post]
func (rc *RoleController) CreateRole(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	//uc := service.NewUserService()
	sort, ctxUser, err := rc.userService.GetCurrentUserMinRoleSort(c)
	if err != nil {
		util.Fail(c, nil, "获取当前用户最高角色等级失败: "+err.Error())
		return
	}
	if sort >= req.Sort {
		util.Fail(c, nil, "不能创建比自己等级高或相同等级的角色")
		return
	}
	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.UserName,
	}
	err = rc.roleService.CreateRole(&role)
	if err != nil {
		util.Fail(c, nil, "创建角色失败: "+err.Error())
		return
	}
	util.Success(c, nil, "创建角色成功")
}

// UpdateRoleById updates an existing role by Id
// @Summary Update role
// @Description Update an existing role by Id
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role Id"
// @Param role body dto.CreateRoleRequest true "Update role request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/role/{roleId} [put]
func (rc *RoleController) UpdateRoleById(c *gin.Context) {
	var req dto.CreateRoleRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		util.Fail(c, nil, "角色ID不正确")
		return
	}
	//ur := service.NewUserService()
	minSort, ctxUser, err := rc.userService.GetCurrentUserMinRoleSort(c)
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	roles, err := rc.roleService.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		util.Fail(c, nil, "未获取到角色信息")
		return
	}
	if minSort >= roles[0].Sort {
		util.Fail(c, nil, "不能更新比自己角色等级高或相等的角色")
		return
	}
	if minSort >= req.Sort {
		util.Fail(c, nil, "不能把角色等级更新得比当前用户的等级高或相同")
		return
	}
	role := model.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.UserName,
	}
	role.Id = uint(roleId)
	err = rc.roleService.UpdateRoleById(&role)
	if err != nil {
		util.Fail(c, nil, "更新角色失败: "+err.Error())
		return
	}
	if req.Keyword != roles[0].Keyword {
		rolePolicies, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if err2 != nil {
			util.Fail(c, nil, "获取角色关键字关联的权限接口失败")
			return
		}
		if len(rolePolicies) == 0 {
			util.Success(c, nil, "更新角色成功")
			return
		}
		rolePoliciesCopy := make([][]string, 0)
		for _, policy := range rolePolicies {
			policyCopy := make([]string, len(policy))
			copy(policyCopy, policy)
			rolePoliciesCopy = append(rolePoliciesCopy, policyCopy)
			policy[0] = req.Keyword
		}
		isAdded, _ := config.CasbinEnforcer.AddPolicies(rolePolicies)
		if !isAdded {
			util.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败")
			return
		}
		isRemoved, _ := config.CasbinEnforcer.RemovePolicies(rolePoliciesCopy)
		if !isRemoved {
			util.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败")
			return
		}
		err := config.CasbinEnforcer.LoadPolicy()
		if err != nil {
			util.Fail(c, nil, "更新角色成功，但角色关键字关联角色的权限接口策略加载失败")
			return
		}
	}
	rc.userService.ClearUserInfoCache()
	util.Success(c, nil, "更新角色成功")
}

// GetRoleMenusById retrieves the menus for a role by Id
// @Summary Get role menus by Id
// @Description Get the menus for a role by Id
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role Id"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/role/{roleId}/menus [get]
func (rc *RoleController) GetRoleMenusById(c *gin.Context) {
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		util.Fail(c, nil, "角色ID不正确")
		return
	}
	menus, err := rc.roleService.GetRoleMenusById(uint(roleId))
	if err != nil {
		util.Fail(c, nil, "获取角色的权限菜单失败: "+err.Error())
		return
	}
	util.Success(c, gin.H{"Menus": menus}, "获取角色的权限菜单成功")
}

// UpdateRoleMenusById updates the menus for a role by Id
// @Summary Update role menus by Id
// @Description Update the menus for a role by Id
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role Id"
// @Param menus body dto.IdListRequest true "Update role menus request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/role/{roleId}/menus [put]
func (rc *RoleController) UpdateRoleMenusById(c *gin.Context) {
	var req dto.IdListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		util.Fail(c, nil, "角色ID不正确")
		return
	}
	// 根据path中的角色Id获取该角色信息
	roles, err := rc.roleService.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		util.Fail(c, nil, "未获取到角色信息")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	//ur := service.NewUserService()
	minSort, ctxUser, err := rc.userService.GetCurrentUserMinRoleSort(c)
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			util.Fail(c, nil, "不能更新比自己角色等级高或相等角色的权限菜单")
			return
		}
	}

	// 获取当前用户所拥有的权限菜单
	mr := service.NewMenuService()
	ctxUserMenus, err := mr.GetUserMenusByUserId(ctxUser.Id)
	if err != nil {
		util.Fail(c, nil, "获取当前用户的可访问菜单列表失败: "+err.Error())
		return
	}

	// 获取当前用户所拥有的权限菜单Id
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.Id)
	}

	// 前端传来最新的MenuIds集合
	menuIds := req.Ids

	// 用户需要修改的菜单集合
	reqMenus := make([]*model.Menu, 0)

	// (非管理员)不能把角色的权限菜单设置的比当前用户所拥有的权限菜单多
	if minSort != 1 {
		for _, id := range menuIds {
			if !slices.Contains(ctxUserMenusIds, id) {
				util.Fail(c, nil, fmt.Sprintf("无权设置ID为%d的菜单", id))
				return
			}
		}

		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.Id {
					reqMenus = append(reqMenus, menu)
					break
				}
			}
		}
	} else {
		// 管理员随意设置
		// 根据menuIds查询查询菜单
		menus, err := mr.GetMenus()
		if err != nil {
			util.Fail(c, nil, "获取菜单列表失败: "+err.Error())
			return
		}
		for _, menuId := range menuIds {
			for _, menu := range menus {
				if menuId == menu.Id {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus

	err = rc.roleService.UpdateRoleMenus(&roles[0])
	if err != nil {
		util.Fail(c, nil, "更新角色的权限菜单失败: "+err.Error())
		return
	}

	util.Success(c, nil, "更新角色的权限菜单成功")

}

// GetRoleApisById retrieves the APIs for a role by Id
// @Summary Get role APIs by Id
// @Description Get the APIs for a role by Id
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role Id"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/role/{roleId}/apis [get]
func (rc *RoleController) GetRoleApisById(c *gin.Context) {
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		util.Fail(c, nil, "角色ID不正确")
		return
	}
	roles, err := rc.roleService.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		util.Fail(c, nil, "未获取到角色信息")
		return
	}
	keyword := roles[0].Keyword
	apis, err := rc.roleService.GetRoleApisByRoleKeyword(keyword)
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	util.Success(c, gin.H{"Apis": apis}, "获取角色的权限接口成功")
}

// UpdateRoleApisById updates the APIs for a role by Id
// @Summary Update role APIs by Id
// @Description Update the APIs for a role by Id
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role Id"
// @Param apis body dto.IdListRequest true "Update role APIs request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/role/{roleId}/apis [put]
func (rc *RoleController) UpdateRoleApisById(c *gin.Context) {
	var req dto.IdListRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		util.Fail(c, nil, errStr)
		return
	}
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		util.Fail(c, nil, "角色ID不正确")
		return
	}
	roles, err := rc.roleService.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		util.Fail(c, nil, "未获取到角色信息")
		return
	}
	//ur := service.NewUserService()
	minSort, ctxUser, err := rc.userService.GetCurrentUserMinRoleSort(c)
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			util.Fail(c, nil, "不能更新比自己角色等级高或相等角色的权限接口")
			return
		}
	}
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		if err2 != nil {
			util.Fail(c, nil, "获取当前用户的角色关键字关联的权限接口失败")
			return
		}
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}
	apiIds := req.Ids
	ar := service.NewApiService()
	apis, err := ar.GetApisById(apiIds)
	if err != nil {
		util.Fail(c, nil, "根据接口ID获取接口信息失败")
		return
	}
	reqRolePolicies := make([][]string, 0)
	for _, api := range apis {
		reqRolePolicies = append(reqRolePolicies, []string{
			roles[0].Keyword, api.Path, api.Method,
		})
	}
	if minSort != 1 {
		for _, reqPolicy := range reqRolePolicies {
			if !funk.Contains(ctxRolesPolicies, reqPolicy) {
				util.Fail(c, nil, fmt.Sprintf("无权设置路径为%s,请求方式为%s的接口", reqPolicy[1], reqPolicy[2]))
				return
			}
		}
	}
	err = rc.roleService.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	util.Success(c, nil, "更新角色的权限接口成功")
}

// BatchDeleteRoleByIds deletes multiple roles by their Ids
// @Summary Batch delete roles
// @Description Delete multiple roles by their Ids
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleIds body dto.IdListRequest true "Delete role request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/role/batch_delete [delete]
func (rc *RoleController) BatchDeleteRoleByIds(c *gin.Context) {
	var req dto.IdListRequest
	if err := c.ShouldBind(&req); err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	minSort, _, err := rc.userService.GetCurrentUserMinRoleSort(c)
	if err != nil {
		util.Fail(c, nil, err.Error())
		return
	}
	roleIds := req.Ids
	roles, err := rc.roleService.GetRolesByIds(roleIds)
	if err != nil {
		util.Fail(c, nil, "获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		util.Fail(c, nil, "未获取到角色信息")
		return
	}
	for _, role := range roles {
		if minSort >= role.Sort {
			util.Fail(c, nil, "不能删除比自己角色等级高或相等的角色")
			return
		}
	}
	err = rc.roleService.BatchDeleteRoleByIds(roleIds)
	if err != nil {
		util.Fail(c, nil, "删除角色失败")
		return
	}
	rc.userService.ClearUserInfoCache()
	util.Success(c, nil, "删除角色成功")
}
