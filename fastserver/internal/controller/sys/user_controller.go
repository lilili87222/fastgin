package sys

import (
	"fastgin/config"
	"fastgin/internal/bean"
	"fastgin/internal/controller"
	sys2 "fastgin/internal/dao/sys"
	"fastgin/internal/model/sys"
	util2 "fastgin/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/thoas/go-funk"
	"strconv"
)

type UserController struct {
	UserRepository sys2.UserRepository
}

// 构造函数
func NewUserController() UserController {
	userRepository := sys2.NewUserRepository()
	userController := UserController{UserRepository: userRepository}
	return userController
}

// 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户信息
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/info [get]
func (uc UserController) GetUserInfo(c *gin.Context) {
	user, err := uc.UserRepository.GetCurrentUser(c)
	if err != nil {
		controller.Fail(c, nil, "获取当前用户信息失败: "+err.Error())
		return
	}
	userInfoDto := ToUserInfoDto(user)
	controller.Success(c, gin.H{
		"userInfo": userInfoDto,
	}, "获取当前用户信息成功")
}

// 获取用户列表
// @Summary 获取用户列表
// @Description 获取用户列表
// @Tags User
// @Accept json
// @Produce json
// @Param UserListRequest body bean.UserListRequest true "User list request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/list [post]
func (uc UserController) GetUsers(c *gin.Context) {
	var req bean.UserListRequest
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

	// 获取
	users, total, err := uc.UserRepository.GetUsers(&req)
	if err != nil {
		controller.Fail(c, nil, "获取用户列表失败: "+err.Error())
		return
	}
	controller.Success(c, gin.H{"users": ToUsersDto(users), "total": total}, "获取用户列表成功")
}

// 更新用户登录密码
// @Summary 更新用户登录密码
// @Description 更新用户登录密码
// @Tags User
// @Accept json
// @Produce json
// @Param ChangePwdRequest body bean.ChangePwdRequest true "Change password request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/change_pwd [post]
func (uc UserController) ChangePwd(c *gin.Context) {
	var req bean.ChangePwdRequest

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

	// 前端传来的密码是rsa加密的,先解密
	// 密码通过RSA解密
	//decodeOldPassword, err := util2.RSADecrypt([]byte(req.OldPassword), config.Conf.System.RSAPrivateBytes)
	//if err != nil {
	//	controller.Fail(c, nil, err.Error())
	//	return
	//}
	//decodeNewPassword, err := util2.RSADecrypt([]byte(req.NewPassword), config.Conf.System.RSAPrivateBytes)
	//if err != nil {
	//	controller.Fail(c, nil, err.Error())
	//	return
	//}
	//req.OldPassword = string(decodeOldPassword)
	//req.NewPassword = string(decodeNewPassword)

	// 获取当前用户
	user, err := uc.UserRepository.GetCurrentUser(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	// 获取用户的真实正确密码
	correctPasswd := user.Password
	// 判断前端请求的密码是否等于真实密码
	err = util2.ComparePasswd(correctPasswd, req.OldPassword)
	if err != nil {
		controller.Fail(c, nil, "原密码有误")
		return
	}
	// 更新密码
	err = uc.UserRepository.ChangePwd(user.Username, util2.GenPasswd(req.NewPassword))
	if err != nil {
		controller.Fail(c, nil, "更新密码失败: "+err.Error())
		return
	}
	controller.Success(c, nil, "更新密码成功")
}

// 创建用户
// @Summary 创建用户
// @Description 创建用户
// @Tags User
// @Accept json
// @Produce json
// @Param CreateUserRequest body bean.CreateUserRequest true "Create user request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/create [post]
func (uc UserController) CreateUser(c *gin.Context) {
	var req bean.CreateUserRequest
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

	// 密码通过RSA解密
	// 密码不为空就解密
	//if req.Password != "" {
	//	decodeData, err := util2.RSADecrypt([]byte(req.Password), config.Conf.System.RSAPrivateBytes)
	//	if err != nil {
	//		controller.Fail(c, nil, err.Error())
	//		return
	//	}
	//	req.Password = string(decodeData)
	//	if len(req.Password) < 6 {
	//		controller.Fail(c, nil, "密码长度至少为6位")
	//		return
	//	}
	//}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	currentRoleSortMin, ctxUser, err := uc.UserRepository.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := sys2.NewRoleRepository()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		controller.Fail(c, nil, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		controller.Fail(c, nil, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := uint(funk.MinInt(reqRoleSorts))

	// 当前用户的角色排序最小值 需要小于 前端传来的角色排序最小值（用户不能创建比自己等级高的或者相同等级的用户）
	if currentRoleSortMin >= reqRoleSortMin {
		controller.Fail(c, nil, "用户不能创建比自己等级高的或者相同等级的用户")
		return
	}

	// 密码为空就默认123456
	if req.Password == "" {
		req.Password = "123456"
	}
	user := sys.User{
		Username:     req.Username,
		Password:     util2.GenPasswd(req.Password),
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}

	err = uc.UserRepository.CreateUser(&user)
	if err != nil {
		controller.Fail(c, nil, "创建用户失败: "+err.Error())
		return
	}
	controller.Success(c, nil, "创建用户成功")

}

// 更新用户
// @Summary 更新用户
// @Description 更新用户
// @Tags User
// @Accept json
// @Produce json
// @Param userId path int true "User ID"
// @Param CreateUserRequest body bean.CreateUserRequest true "Update user request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/update/{userId} [put]
func (uc UserController) UpdateUserById(c *gin.Context) {
	var req bean.CreateUserRequest
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

	//获取path中的userId
	userId, _ := strconv.Atoi(c.Param("userId"))
	if userId <= 0 {
		controller.Fail(c, nil, "用户ID不正确")
		return
	}

	// 根据path中的userId获取用户信息
	oldUser, err := uc.UserRepository.GetUserById(uint(userId))
	if err != nil {
		controller.Fail(c, nil, "获取需要更新的用户信息失败: "+err.Error())
		return
	}

	// 获取当前用户
	ctxUser, err := uc.UserRepository.GetCurrentUser(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
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
		currentRoleIds = append(currentRoleIds, role.ID)
	}
	// 当前用户角色排序最小值（最高等级角色）
	currentRoleSortMin := funk.MinInt(currentRoleSorts)

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := sys2.NewRoleRepository()
	roles, err := rr.GetRolesByIds(reqRoleIds)
	if err != nil {
		controller.Fail(c, nil, "根据角色ID获取角色信息失败: "+err.Error())
		return
	}
	if len(roles) == 0 {
		controller.Fail(c, nil, "未获取到角色信息")
		return
	}
	var reqRoleSorts []int
	for _, role := range roles {
		reqRoleSorts = append(reqRoleSorts, int(role.Sort))
	}
	// 前端传来用户角色排序最小值（最高等级角色）
	reqRoleSortMin := funk.MinInt(reqRoleSorts)

	user := sys.User{
		Model:        oldUser.Model,
		Username:     req.Username,
		Password:     oldUser.Password,
		Mobile:       req.Mobile,
		Avatar:       req.Avatar,
		Nickname:     &req.Nickname,
		Introduction: &req.Introduction,
		Status:       req.Status,
		Creator:      ctxUser.Username,
		Roles:        roles,
	}
	// 判断是更新自己还是更新别人
	if userId == int(ctxUser.ID) {
		// 如果是更新自己
		// 不能禁用自己
		if req.Status == 2 {
			controller.Fail(c, nil, "不能禁用自己")
			return
		}
		// 不能更改自己的角色
		reqDiff, currentDiff := funk.Difference(req.RoleIds, currentRoleIds)
		if len(reqDiff.([]uint)) > 0 || len(currentDiff.([]uint)) > 0 {
			controller.Fail(c, nil, "不能更改自己的角色")
			return
		}

		// 不能更新自己的密码，只能在个人中心更新
		if req.Password != "" {
			controller.Fail(c, nil, "请到个人中心更新自身密码")
			return
		}

		// 密码赋值
		user.Password = ctxUser.Password

	} else {
		// 如果是更新别人
		// 用户不能更新比自己角色等级高的或者相同等级的用户
		// 根据path中的userIdID获取用户角色排序最小值
		minRoleSorts, err := uc.UserRepository.GetUserMinRoleSortsByIds([]uint{uint(userId)})
		if err != nil || len(minRoleSorts) == 0 {
			controller.Fail(c, nil, "根据用户ID获取用户角色排序最小值失败")
			return
		}
		if currentRoleSortMin >= minRoleSorts[0] {
			controller.Fail(c, nil, "用户不能更新比自己角色等级高的或者相同等级的用户")
			return
		}

		// 用户不能把别的用户角色等级更新得比自己高或相等
		if currentRoleSortMin >= reqRoleSortMin {
			controller.Fail(c, nil, "用户不能把别的用户角色等级更新得比自己高或相等")
			return
		}

		// 密码赋值
		//if req.Password != "" {
		//	// 密码通过RSA解密
		//	decodeData, err := util2.RSADecrypt([]byte(req.Password), config.Conf.System.RSAPrivateBytes)
		//	if err != nil {
		//		controller.Fail(c, nil, err.Error())
		//		return
		//	}
		//	req.Password = string(decodeData)
		//	user.Password = util2.GenPasswd(req.Password)
		//}
		user.Password = util2.GenPasswd(req.Password)

	}

	// 更新用户
	err = uc.UserRepository.UpdateUser(&user)
	if err != nil {
		controller.Fail(c, nil, "更新用户失败: "+err.Error())
		return
	}
	controller.Success(c, nil, "更新用户成功")

}

// 批量删除用户
// @Summary 批量删除用户
// @Description 批量删除用户
// @Tags User
// @Accept json
// @Produce json
// @Param DeleteUserRequest body bean.DeleteUserRequest true "Delete user request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/batch_delete [delete]
func (uc UserController) BatchDeleteUserByIds(c *gin.Context) {
	var req bean.DeleteUserRequest
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

	// 前端传来的用户ID
	reqUserIds := req.UserIds
	// 根据用户ID获取用户角色排序最小值
	roleMinSortList, err := uc.UserRepository.GetUserMinRoleSortsByIds(reqUserIds)
	if err != nil || len(roleMinSortList) == 0 {
		controller.Fail(c, nil, "根据用户ID获取用户角色排序最小值失败")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := uc.UserRepository.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	currentRoleSortMin := int(minSort)

	// 不能删除自己
	if funk.Contains(reqUserIds, ctxUser.ID) {
		controller.Fail(c, nil, "用户不能删除自己")
		return
	}

	// 不能删除比自己角色排序低(等级高)的用户
	for _, sort := range roleMinSortList {
		if currentRoleSortMin >= sort {
			controller.Fail(c, nil, "用户不能删除比自己角色等级高的用户")
			return
		}
	}

	err = uc.UserRepository.BatchDeleteUserByIds(reqUserIds)
	if err != nil {
		controller.Fail(c, nil, "删除用户失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "删除用户成功")

}