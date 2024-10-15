package sys

import (
	"fastgin/config"
	"fastgin/internal/controller"
	sys2 "fastgin/internal/dao/sys"
	"fastgin/internal/model/sys"
	"fastgin/internal/model/sys/request"
	sys3 "fastgin/internal/service/sys"
	util2 "fastgin/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"github.com/thoas/go-funk"
	"slices"
	"strconv"
)

type UserController struct {
	UserDao sys2.UserDao
}

// 构造函数
func NewUserController() UserController {
	userDao := sys2.NewUserDao()
	userController := UserController{UserDao: userDao}
	return userController
}

// 获取当前登录用户信息
// @Summary 获取当前登录用户信息
// @Description 获取当前登录用户信息
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/info [get]
func (uc UserController) GetUserInfo(c *gin.Context) {
	type UserInfoDto struct {
		ID           uint        `json:"id"`
		Username     string      `json:"username"`
		Mobile       string      `json:"mobile"`
		Avatar       string      `json:"avatar"`
		Nickname     string      `json:"nickname"`
		Introduction string      `json:"introduction"`
		Roles        []*sys.Role `json:"roles"`
	}
	user, err := uc.UserDao.GetCurrentUser(c)
	if err != nil {
		controller.Fail(c, nil, "获取当前用户信息失败: "+err.Error())
		return
	}
	userInfoDto := UserInfoDto{}
	copier.Copy(&userInfoDto, &user)
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
// @Param Authorization header string true "Bearer token"
// @Param UserListRequest body request.UserListRequest true "User list request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/list [post]
func (uc UserController) GetUsers(c *gin.Context) {

	var req request.UserListRequest
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
	users, total, err := uc.UserDao.GetUsers(&req)
	if err != nil {
		controller.Fail(c, nil, "获取用户列表失败: "+err.Error())
		return
	}

	type UsersDto struct {
		ID           uint   `json:"ID"`
		Username     string `json:"username"`
		Mobile       string `json:"mobile"`
		Avatar       string `json:"avatar"`
		Nickname     string `json:"nickname"`
		Introduction string `json:"introduction"`
		Status       uint   `json:"status"`
		Creator      string `json:"creator"`
		RoleIds      []uint `json:"roleIds"`
	}
	ToUsersDto := func(userList []*sys.User) []UsersDto {
		var users []UsersDto
		for _, user := range userList {
			userDto := UsersDto{}
			copier.Copy(&userDto, user)
			userDto.RoleIds = user.GetRoleIds()
			users = append(users, userDto)
		}
		return users
	}
	controller.Success(c, gin.H{"users": ToUsersDto(users), "total": total}, "获取用户列表成功")
}

// 更新用户登录密码
// @Summary 更新用户登录密码
// @Description 更新用户登录密码
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param ChangePwdRequest body request.ChangePwdRequest true "Change password request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/change_pwd [post]
func (uc UserController) ChangePwd(c *gin.Context) {
	var req request.ChangePwdRequest

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
	//decodeOldPassword, err := util2.RSADecrypt([]byte(req.OldPassword), config.Instance.System.RSAPrivateBytes)
	//if err != nil {
	//	controller.Fail(c, nil, err.Error())
	//	return
	//}
	//decodeNewPassword, err := util2.RSADecrypt([]byte(req.NewPassword), config.Instance.System.RSAPrivateBytes)
	//if err != nil {
	//	controller.Fail(c, nil, err.Error())
	//	return
	//}
	//req.OldPassword = string(decodeOldPassword)
	//req.NewPassword = string(decodeNewPassword)

	// 获取当前用户
	user, err := uc.UserDao.GetCurrentUser(c)
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
	err = uc.UserDao.ChangePwd(user.Username, util2.GenPasswd(req.NewPassword))
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
// @Param Authorization header string true "Bearer token"
// @Param CreateUserRequest body request.CreateUserRequest true "Create user request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/create [post]
func (uc UserController) CreateUser(c *gin.Context) {
	var req request.CreateUserRequest
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
	//	decodeData, err := util2.RSADecrypt([]byte(req.Password), config.Instance.System.RSAPrivateBytes)
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
	currentRoleSortMin, ctxUser, err := uc.UserDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := sys3.NewRoleDao()
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
	//reqRoleSortMin := uint(funk.MinInt(reqRoleSorts))
	reqRoleSortMin := uint(slices.Min(reqRoleSorts))

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

	err = uc.UserDao.CreateUser(&user)
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
// @Param Authorization header string true "Bearer token"
// @Param userId path int true "User ID"
// @Param CreateUserRequest body request.CreateUserRequest true "Update user request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/update/{userId} [put]
func (uc UserController) UpdateUserById(c *gin.Context) {
	var req request.CreateUserRequest
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
	oldUser, err := uc.UserDao.GetUserById(uint(userId))
	if err != nil {
		controller.Fail(c, nil, "获取需要更新的用户信息失败: "+err.Error())
		return
	}

	// 获取当前用户
	ctxUser, err := uc.UserDao.GetCurrentUser(c)
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
	currentRoleSortMin := slices.Min(currentRoleSorts)

	// 获取前端传来的用户角色id
	reqRoleIds := req.RoleIds
	// 根据角色id获取角色
	rr := sys3.NewRoleDao()
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
	reqRoleSortMin := slices.Min(reqRoleSorts)

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
		minRoleSorts, err := uc.UserDao.GetUserMinRoleSortsByIds([]uint{uint(userId)})
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
		//	decodeData, err := util2.RSADecrypt([]byte(req.Password), config.Instance.System.RSAPrivateBytes)
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
	err = uc.UserDao.UpdateUser(&user)
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
// @Param Authorization header string true "Bearer token"
// @Param DeleteUserRequest body request.DeleteUserRequest true "Delete user request"
// @Success 200 {object} controller.ResponseBody
// @Failure 400 {object} controller.ResponseBody
// @Router /api/auth/user/batch_delete [delete]
func (uc UserController) BatchDeleteUserByIds(c *gin.Context) {
	var req request.DeleteUserRequest
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
	roleMinSortList, err := uc.UserDao.GetUserMinRoleSortsByIds(reqUserIds)
	if err != nil || len(roleMinSortList) == 0 {
		controller.Fail(c, nil, "根据用户ID获取用户角色排序最小值失败")
		return
	}

	// 当前用户角色排序最小值（最高等级角色）以及当前用户
	minSort, ctxUser, err := uc.UserDao.GetCurrentUserMinRoleSort(c)
	if err != nil {
		controller.Fail(c, nil, err.Error())
		return
	}
	currentRoleSortMin := int(minSort)

	// 不能删除自己
	if slices.Contains(reqUserIds, ctxUser.ID) {
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

	err = uc.UserDao.BatchDeleteUserByIds(reqUserIds)
	if err != nil {
		controller.Fail(c, nil, "删除用户失败: "+err.Error())
		return
	}

	controller.Success(c, nil, "删除用户成功")

}
