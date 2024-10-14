package sys

import (
	"fastgin/config"
	"fastgin/internal/bean"
	"fastgin/internal/controller"
	sys2 "fastgin/internal/dao/sys"
	"fastgin/internal/model/sys"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
	"slices"
	"strconv"
)

// RoleController handles role-related requests
type RoleController struct {
	RoleDao sys2.RoleDao
}

// NewRoleController creates a new RoleController
func NewRoleController() RoleController {
	roleDao := sys2.NewRoleDao()
	roleController := RoleController{RoleDao: roleDao}
	return roleController
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
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/roles [get]
func (rc RoleController) GetRoles(c *gin.Context) {
	var req bean.RoleListRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	roles, total, err := rc.RoleDao.GetRoles(&req)
	if err != nil {
		controller.Fail(c, nil, "获取角色列表失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"roles": roles, "total": total}, "获取角色列表成功")
}

// CreateRole creates a new role
// @Summary Create role
// @Description Create a new role
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param role body bean.CreateRoleRequest true "Create role request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/role [post]
func (rc RoleController) CreateRole(c *gin.Context) {
	var req bean.CreateRoleRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	uc := sys2.NewUserDao()
	sort, ctxUser, err := uc.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, "获取当前用户最高角色等级失败: "+err.Error())
		return
	}
	if sort >= req.Sort {
		controller.Fail(c, nil, "不能创建比自己等级高或相同等级的角色")
		return
	}
	role := sys.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}
	err = rc.RoleDao.CreateRole(&role)
	if err != nil {
		controller.Fail(c, nil, "创建角色失败: "+err.Error())
		return
	}
	controller.Success(c, nil, "创建角色成功")
}

// UpdateRoleById updates an existing role by ID
// @Summary Update role
// @Description Update an existing role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role ID"
// @Param role body bean.CreateRoleRequest true "Update role request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/role/{roleId} [put]
func (rc RoleController) UpdateRoleById(c *gin.Context) {
	var req bean.CreateRoleRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		controller.Fail(c, nil, "角色ID不正确")
		return
	}
	ur := sys2.NewUserDao()
	minSort, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	roles, err := rc.RoleDao.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		controller.Fail(c, nil, "未获取到角色信息")
		return
	}
	if minSort >= roles[0].Sort {
		controller.Fail(c, nil, "不能更新比自己角色等级高或相等的角色")
		return
	}
	if minSort >= req.Sort {
		controller.Fail(c, nil, "不能把角色等级更新得比当前用户的等级高或相同")
		return
	}
	role := sys.Role{
		Name:    req.Name,
		Keyword: req.Keyword,
		Desc:    &req.Desc,
		Status:  req.Status,
		Sort:    req.Sort,
		Creator: ctxUser.Username,
	}
	err = rc.RoleDao.UpdateRoleById(uint(roleId), &role)
	if err != nil {
		controller.Fail(c, nil, "更新角色失败: "+err.Error())
		return
	}
	if req.Keyword != roles[0].Keyword {
		rolePolicies, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, roles[0].Keyword)
		if err2 != nil {
			controller.Fail(c, nil, "获取角色关键字关联的权限接口失败")
			return
		}
		if len(rolePolicies) == 0 {
			controller.Success(c, nil, "更新角色成功")
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
			controller.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败")
			return
		}
		isRemoved, _ := config.CasbinEnforcer.RemovePolicies(rolePoliciesCopy)
		if !isRemoved {
			controller.Fail(c, nil, "更新角色成功，但角色关键字关联的权限接口更新失败")
			return
		}
		err := config.CasbinEnforcer.LoadPolicy()
		if err != nil {
			controller.Fail(c, nil, "更新角色成功，但角色关键字关联角色的权限接口策略加载失败")
			return
		}
	}
	ur.ClearUserInfoCache()
	controller.Success(c, nil, "更新角色成功")
}

// GetRoleMenusById retrieves the menus for a role by ID
// @Summary Get role menus by ID
// @Description Get the menus for a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role ID"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/role/{roleId}/menus [get]
func (rc RoleController) GetRoleMenusById(c *gin.Context) {
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		controller.Fail(c, nil, "角色ID不正确")
		return
	}
	menus, err := rc.RoleDao.GetRoleMenusById(uint(roleId))
	if err != nil {
		controller.Fail(c, nil, "获取角色的权限菜单失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"menus": menus}, "获取角色的权限菜单成功")
}

// UpdateRoleMenusById updates the menus for a role by ID
// @Summary Update role menus by ID
// @Description Update the menus for a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role ID"
// @Param menus body bean.UpdateRoleMenusRequest true "Update role menus request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/role/{roleId}/menus [put]
func (rc RoleController) UpdateRoleMenusById(c *gin.Context) {
	var req bean.UpdateRoleMenusRequest
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
	// 获取path中的roleId
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		controller.Fail(c, nil, "角色ID不正确")
		return
	}
	// 根据path中的角色ID获取该角色信息
	roles, err := rc.RoleDao.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		controller.Fail(c, nil, "未获取到角色信息")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	ur := sys2.NewUserDao()
	minSort, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}

	// (非管理员)不能更新比自己角色等级高或相等角色的权限菜单
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			controller.Fail(c, nil, "不能更新比自己角色等级高或相等角色的权限菜单")
			return
		}
	}

	// 获取当前用户所拥有的权限菜单
	mr := sys2.NewMenuDao()
	ctxUserMenus, err := mr.GetUserMenusByUserId(ctxUser.ID)
	if err != nil {
		controller.Fail(c, nil, "获取当前用户的可访问菜单列表失败: "+err.Error())
		return
	}

	// 获取当前用户所拥有的权限菜单ID
	ctxUserMenusIds := make([]uint, 0)
	for _, menu := range ctxUserMenus {
		ctxUserMenusIds = append(ctxUserMenusIds, menu.ID)
	}

	// 前端传来最新的MenuIds集合
	menuIds := req.MenuIds

	// 用户需要修改的菜单集合
	reqMenus := make([]*sys.Menu, 0)

	// (非管理员)不能把角色的权限菜单设置的比当前用户所拥有的权限菜单多
	if minSort != 1 {
		for _, id := range menuIds {
			if !slices.Contains(ctxUserMenusIds, id) {
				controller.Fail(c, nil, fmt.Sprintf("无权设置ID为%d的菜单", id))
				return
			}
		}

		for _, id := range menuIds {
			for _, menu := range ctxUserMenus {
				if id == menu.ID {
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
			controller.Fail(c, nil, "获取菜单列表失败: "+err.Error())
			return
		}
		for _, menuId := range menuIds {
			for _, menu := range menus {
				if menuId == menu.ID {
					reqMenus = append(reqMenus, menu)
				}
			}
		}
	}

	roles[0].Menus = reqMenus

	err = rc.RoleDao.UpdateRoleMenus(roles[0])
	if err != nil {
		controller.Fail(c, nil, "更新角色的权限菜单失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "更新角色的权限菜单成功")

}

// GetRoleApisById retrieves the APIs for a role by ID
// @Summary Get role APIs by ID
// @Description Get the APIs for a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role ID"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/role/{roleId}/apis [get]
func (rc RoleController) GetRoleApisById(c *gin.Context) {
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		controller.Fail(c, nil, "角色ID不正确")
		return
	}
	roles, err := rc.RoleDao.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		controller.Fail(c, nil, "未获取到角色信息")
		return
	}
	keyword := roles[0].Keyword
	apis, err := rc.RoleDao.GetRoleApisByRoleKeyword(keyword)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	controller.Success(c, gin.H{"apis": apis}, "获取角色的权限接口成功")
}

// UpdateRoleApisById updates the APIs for a role by ID
// @Summary Update role APIs by ID
// @Description Update the APIs for a role by ID
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleId path int true "Role ID"
// @Param apis body bean.UpdateRoleApisRequest true "Update role APIs request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/role/{roleId}/apis [put]
func (rc RoleController) UpdateRoleApisById(c *gin.Context) {
	var req bean.UpdateRoleApisRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	roleId, _ := strconv.Atoi(c.Param("roleId"))
	if roleId <= 0 {
		controller.Fail(c, nil, "角色ID不正确")
		return
	}
	roles, err := rc.RoleDao.GetRolesByIds([]uint{uint(roleId)})
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if len(roles) == 0 {
		controller.Fail(c, nil, "未获取到角色信息")
		return
	}
	ur := sys2.NewUserDao()
	minSort, ctxUser, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if minSort != 1 {
		if minSort >= roles[0].Sort {
			controller.Fail(c, nil, "不能更新比自己角色等级高或相等角色的权限接口")
			return
		}
	}
	ctxRoles := ctxUser.Roles
	ctxRolesPolicies := make([][]string, 0)
	for _, role := range ctxRoles {
		policy, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, role.Keyword)
		if err2 != nil {
			controller.Fail(c, nil, "获取当前用户的角色关键字关联的权限接口失败")
			return
		}
		ctxRolesPolicies = append(ctxRolesPolicies, policy...)
	}
	for _, policy := range ctxRolesPolicies {
		policy[0] = roles[0].Keyword
	}
	apiIds := req.ApiIds
	ar := sys2.NewApiDao()
	apis, err := ar.GetApisById(apiIds)
	if err != nil {
		controller.Fail(c, nil, "根据接口ID获取接口信息失败")
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
				controller.Fail(c, nil, fmt.Sprintf("无权设置路径为%s,请求方式为%s的接口", reqPolicy[1], reqPolicy[2]))
				return
			}
		}
	}
	err = rc.RoleDao.UpdateRoleApis(roles[0].Keyword, reqRolePolicies)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	controller.Success(c, nil, "更新角色的权限接口成功")
}

// BatchDeleteRoleByIds deletes multiple roles by their IDs
// @Summary Batch delete roles
// @Description Delete multiple roles by their IDs
// @Tags Role
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param roleIds body bean.DeleteRoleRequest true "Delete role request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/role/batch_delete [delete]
func (rc RoleController) BatchDeleteRoleByIds(c *gin.Context) {
	var req bean.DeleteRoleRequest
	if err := c.ShouldBind(&req); err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		controller.Fail(c, nil, errStr)
		return
	}
	ur := sys2.NewUserDao()
	minSort, _, err := ur.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	roleIds := req.RoleIds
	roles, err := rc.RoleDao.GetRolesByIds(roleIds)
	if err != nil {
		controller.Fail(c, nil, "获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		controller.Fail(c, nil, "未获取到角色信息")
		return
	}
	for _, role := range roles {
		if minSort >= role.Sort {
			controller.Fail(c, nil, "不能删除比自己角色等级高或相等的角色")
			return
		}
	}
	err = rc.RoleDao.BatchDeleteRoleByIds(roleIds)
	if err != nil {
		controller.Fail(c, nil, "删除角色失败")
		return
	}
	ur.ClearUserInfoCache()
	controller.Success(c, nil, "删除角色成功")
}
