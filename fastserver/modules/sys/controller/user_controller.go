package controller

import (
	httpz2 "fastgin/common/httpz"
	util2 "fastgin/common/util"
	"fastgin/config"
	"fastgin/modules/sys/dto"
	"fastgin/modules/sys/model"
	"fastgin/modules/sys/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
	"slices"
	"strconv"
)

type UserController struct {
	userService *service.UserService
}

// 构造函数
func NewUserController() *UserController {
	return &UserController{userService: service.NewUserService()}
}

// 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户信息
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/user/info [get]
func (uc *UserController) GetUserInfo(c *gin.Context) {
	user, err := uc.userService.GetCurrentUser(c)
	if err != nil {
		httpz2.ServerError(c, "获取当前用户信息失败: "+err.Error())
		return
	}
	userMap, e := util2.StructToMap(user, true, "Id", "UserName", "Mobile", "Avatar", "NickName", "Introduction", "Roles")
	if e != nil {
		fmt.Println(e.Error())
	}
	//fmt.Println(util.Struct2Json(userMap))
	//fmt.Println(util.Struct2Json(user))
	//userInfoDto := dto.UserInfoDto{}
	//copier.Copy(&userInfoDto, &user)
	httpz2.Success(c, userMap)
}

// 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param UserListRequest body dto.UserListRequest true "User list request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/user/index [get]
func (uc *UserController) GetUsers(c *gin.Context) {
	params, e := httpz2.GetFormData(c)
	if e != nil {
		httpz2.BadRequest(c, e.Error())
		return
	}
	data, total, err := uc.userService.GetUsersWithRoleIds(httpz2.NewSearchRequest(params))
	if err != nil {
		httpz2.ServerError(c, "获取角色列表失败: "+err.Error())
		return
	}
	httpz2.Success(c, gin.H{"Data": data, "Total": total})

	//
	//
	//var req dto.UserListRequest
	//// 参数绑定
	//if err := c.ShouldBind(&req); err != nil {
	//	util.ServerError(c,  err.Error())
	//	return
	//}
	//// 参数校验
	//if err := config.Validate.Struct(&req); err != nil {
	//	errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
	//	util.ServerError(c,  errStr)
	//	return
	//}
	//
	//// 获取
	//users, total, err := uc.userService.GetUsersWithRoleIds(req)
	//if err != nil {
	//	util.ServerError(c,  "获取用户列表失败: "+err.Error())
	//	return
	//}
	//
	//util.Success(c, gin.H{"Users": users, "Total": total}, "获取用户列表成功")
}

// 更新用户登录密码
// @Summary 更新用户登录密码
// @Description 更新用户登录密码
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param ChangePwdRequest body dto.ChangePwdRequest true "Change password request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/user/changePwd [post]
func (uc *UserController) ChangePwd(c *gin.Context) {
	var req dto.ChangePwdRequest

	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		httpz2.BadRequest(c, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		httpz2.BadRequest(c, errStr)
		return
	}

	// 前端传来的密码是rsa加密的,先解密
	// 密码通过RSA解密
	//decodeOldPassword, err := util.RSADecrypt([]byte(req.OldPassword), config.Instance.System.RSAPrivateBytes)
	//if err != nil {
	//	controller.ServerError(c,  err.Error())
	//	return
	//}
	//decodeNewPassword, err := util.RSADecrypt([]byte(req.NewPassword), config.Instance.System.RSAPrivateBytes)
	//if err != nil {
	//	controller.ServerError(c,  err.Error())
	//	return
	//}
	//req.OldPassword = string(decodeOldPassword)
	//req.NewPassword = string(decodeNewPassword)

	// 获取当前用户
	user, err := uc.userService.GetCurrentUser(c)
	if err != nil {
		httpz2.ServerError(c, err.Error())
		return
	}
	// 获取用户的真实正确密码
	correctPasswd := user.Password
	// 判断前端请求的密码是否等于真实密码
	err = util2.ComparePasswd(correctPasswd, req.OldPassword)
	if err != nil {
		httpz2.ServerError(c, "原密码有误")
		return
	}
	// 更新密码
	err = uc.userService.ChangePwd(user.UserName, util2.GenPasswd(req.NewPassword))
	if err != nil {
		httpz2.ServerError(c, "更新密码失败: "+err.Error())
		return
	}
	httpz2.Success(c, nil)
}

// 创建用户
// @Summary 创建用户
// @Description 创建用户
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param CreateUserRequest body dto.CreateUserRequest true "Create user request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/user/index [post]
func (uc *UserController) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		httpz2.BadRequest(c, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		httpz2.BadRequest(c, errStr)
		return
	}

	// 密码通过RSA解密
	// 密码不为空就解密
	//if req.Password != "" {
	//	decodeData, err := util.RSADecrypt([]byte(req.Password), config.Instance.System.RSAPrivateBytes)
	//	if err != nil {
	//		controller.ServerError(c,  err.Error())
	//		return
	//	}
	//	req.Password = string(decodeData)
	//	if len(req.Password) < 6 {
	//		controller.ServerError(c,  "密码长度至少为6位")
	//		return
	//	}
	//}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	currentRoleSortMin, ctxUser, err := uc.userService.GetCurrentUserMinRoleSort(c)
	if err != nil {
		httpz2.ServerError(c, err.Error())
		return
	}

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := service.NewRoleService()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		httpz2.ServerError(c, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		httpz2.ServerError(c, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	//reqRoleSortMin := uint(funk.MinInt(reqRoleSorts))
	reqRoleSortMin := uint(slices.Min(reqRoleSorts))

	// 当前用户的角色排序最小值 需要小于 前端传来的角色排序最小值（用户不能创建比自己等级高的或者相同等级的用户）
	if currentRoleSortMin >= reqRoleSortMin {
		httpz2.ServerError(c, "用户不能创建比自己等级高的或者相同等级的用户")
		return
	}

	// 密码为空就默认123456
	if req.Password == "" {
		req.Password = "123456"
	}
	user := model.User{
		UserName:     req.Username,
		Password:     util2.GenPasswd(req.Password),
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		NickName:     &req.NickName,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.UserName,
		Roles:        roles,
	}

	err = uc.userService.CreateUser(&user)
	if err != nil {
		httpz2.ServerError(c, "创建用户失败: "+err.Error())
		return
	}
	httpz2.Success(c, nil)

}

// 更新用户
// @Summary 更新用户
// @Description 更新用户
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param userId path int true "User Id"
// @Param CreateUserRequest body dto.CreateUserRequest true "Update user request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/user/index/{userId} [put]
func (uc *UserController) Update(c *gin.Context) {
	var req dto.CreateUserRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		httpz2.BadRequest(c, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		httpz2.BadRequest(c, errStr)
		return
	}

	//获取path中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		httpz2.BadRequest(c, "用户ID不正确")
		return
	}

	// 根据path中的userId获取用户信息
	oldUser, err := uc.userService.GetUserById(uint(userId))
	if err != nil {
		httpz2.ServerError(c, "获取需要更新的用户信息失败: "+err.Error())
		return
	}

	// 获取当前用户
	ctxUser, err := uc.userService.GetCurrentUser(c)
	if err != nil {
		httpz2.ServerError(c, err.Error())
		return
	}
	// 获取当前用户的所有角色
	currentRoles := ctxUser.Roles
	// 获取当前用户角色的排序，和前端传来的角色排序做比较
	var currentRoleSorts []int
	// 当前用户角色ID集合
	var currentRoleIds []uint
	for _, role := range currentRoles {
		currentRoleSorts = append(currentRoleSorts, int(role.Sort))
		currentRoleIds = append(currentRoleIds, role.Id)
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := slices.Min(currentRoleSorts)

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := service.NewRoleService()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		httpz2.ServerError(c, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		httpz2.ServerError(c, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := slices.Min(reqRoleSorts)

	user := model.User{
		Model:        oldUser.Model,
		UserName:     req.Username,
		Password:     oldUser.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		NickName:     &req.NickName,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.UserName,
		Roles:        roles,
	}
	// 判断是更新自己还是更新别人
	if userId == int(ctxUser.Id) {
		// 如果是更新自己
		// 不能禁用自己
		if req.Status == 2 {
			httpz2.ServerError(c, "不能禁用自己")
			return
		}
		// 不能更改自己的角色
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]uint)) > 0 || len(currentDiff.([]uint)) > 0 {
			httpz2.ServerError(c, "不能更改自己的角色")
			return
		}

		// 不能更新自己的密码，只能在个人中心更新
		if req.Password != "" {
			httpz2.ServerError(c, "请到个人中心更新自身密码")
			return
		}

		// 密码赋值
		user.Password = ctxUser.Password

	} else {
		// 如果是更新别人
		// 用户不能更新比自己角色等级高的或者相同等级的用户
		// 根据path中的userIdID获取用户角色排序最小值
		minRoleSorts, err := uc.userService.GetUserMinRoleSortsByIds([]uint{uint(userId)})
		if err != nil || len(minRoleSorts) == 0 {
			httpz2.ServerError(c, "根据用户ID获取用户角色排序最小值失败")
			return
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			httpz2.ServerError(c, "用户不能更新比自己角色等级高的或者相同等级的用户")
			return
		}

		// 用户不能把别的用户角色等级更新得比自己高或相等
		if currentRoleSortMin >= reqRoleSortMin {
			httpz2.ServerError(c, "用户不能把别的用户角色等级更新得比自己高或相等")
			return
		}

		// 密码赋值
		//if req.Password != "" {
		//	// 密码通过RSA解密
		//	decodeData, err := util.RSADecrypt([]byte(req.Password), config.Instance.System.RSAPrivateBytes)
		//	if err != nil {
		//		controller.ServerError(c,  err.Error())
		//		return
		//	}
		//	req.Password = string(decodeData)
		//	user.Password = util.GenPasswd(req.Password)
		//}
		user.Password = util2.GenPasswd(req.Password)

	}

	// 更新用户
	err = uc.userService.UpdateUser(&user)
	if err != nil {
		httpz2.ServerError(c, "更新用户失败: "+err.Error())
		return
	}
	httpz2.Success(c, nil)

}

// 批量删除用户
// @Summary 批量删除用户
// @Description 批量删除用户
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param DeleteUserRequest body dto.IdListRequest true "BatchDelete user request"
// @Success 200 {object} util.ResponseBody
// @Failure 400 {object} util.ResponseBody
// @Router /api/auth/user/index [delete]
func (uc *UserController) BatchDeleteUserByIds(c *gin.Context) {
	var req httpz2.IdListRequest
	// 参数绑定
	if err := c.ShouldBind(&req); err != nil {
		httpz2.BadRequest(c, err.Error())
		return
	}
	// 参数校验
	if err := config.Validate.Struct(&req); err != nil {
		errStr := err.(validator.ValidationErrors)[0].Translate(config.Trans)
		httpz2.BadRequest(c, errStr)
		return
	}

	// 前端传来的用户ID
	reqUserIds := req.Ids
	// 根据用户ID获取用户角色排序最小值
	roleMinSortList, err := uc.userService.GetUserMinRoleSortsByIds(reqUserIds)
	if err != nil || len(roleMinSortList) == 0 {
		httpz2.ServerError(c, "根据用户ID获取用户角色排序最小值失败")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := uc.userService.GetCurrentUserMinRoleSort(c)
	if err != nil {
		httpz2.ServerError(c, err.Error())
		return
	}
	currentRoleSortMin := int(minSort)

	// 不能删除自己
	if slices.Contains(reqUserIds, ctxUser.Id) {
		httpz2.ServerError(c, "用户不能删除自己")
		return
	}

	// 不能删除比自己角色排序低(等级高)的用户
	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			httpz2.ServerError(c, "用户不能删除比自己角色等级高的用户")
			return
		}
	}

	err = uc.userService.BatchDeleteUserByIds(reqUserIds)
	if err != nil {
		httpz2.ServerError(c, "删除用户失败: "+err.Error())
		return
	}

	httpz2.Success(c, nil)

}