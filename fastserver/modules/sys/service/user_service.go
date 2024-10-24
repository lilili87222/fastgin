package service

import (
	"errors"
	"fastgin/common/cache"
	"fastgin/common/httpz"
	"fastgin/common/util"
	"fastgin/database"
	"fastgin/modules/sys/dao"
	"fastgin/modules/sys/dto"
	"fastgin/modules/sys/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"slices"
)

type UserService struct {
	userDao *dao.UserDao
}

func NewUserService() *UserService {
	return &UserService{userDao: dao.NewUserDao()}
}

// 登录
func (us *UserService) Login(user *model.User) (*model.User, error) {
	firstUser := us.userDao.GetUserByUsername(user.UserName)
	if firstUser == nil {
		return nil, errors.New("用户不存在")
	}

	if firstUser.Status != 1 {
		return nil, errors.New("用户被禁用")
	}

	isValidate := false
	for _, role := range firstUser.Roles {
		if role.Status == 1 {
			isValidate = true
			break
		}
	}

	if !isValidate {
		return nil, errors.New("用户角色被禁用")
	}

	err := util.ComparePasswd(firstUser.Password, user.Password)
	if err != nil {
		return firstUser, errors.New("密码错误")
	}
	return firstUser, nil
}

// 获取当前登录用户信息
func GetCurrentUser(c *gin.Context) *model.User {
	u, _ := c.Get("user")
	return u.(*model.User)
}

// 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
func GetCurrentUserMinRoleSort(c *gin.Context) (int32, *model.User, error) {
	ctxUser := GetCurrentUser(c)

	currentRoleSorts := make([]int, len(ctxUser.Roles))
	for i, role := range ctxUser.Roles {
		currentRoleSorts[i] = int(role.Sort)
	}

	currentRoleSortMin := slices.Min(currentRoleSorts)
	return int32(currentRoleSortMin), ctxUser, nil
}

// 获取单个用户
func (us *UserService) GetUserById(id uint64) (*model.User, error) {
	return us.userDao.GetUserWithRoles(id)
}

// 获取用户列表
func (us *UserService) GetUsers(req *httpz.SearchRequest) ([]*model.User, int64, error) {
	return database.SearchTable[*model.User](req)
	//return us.userDao.GetUsers(req)
}
func (us *UserService) GetUsersWithRoleIds(req *httpz.SearchRequest) ([]dto.UsersDto, int64, error) {
	userList, i, err := us.GetUsers(req)
	var ids []uint64
	for _, user := range userList {
		ids = append(ids, user.ID)
	}
	userRoleList, e := us.userDao.GetUsersWithRoles(ids)
	if e != nil {
		return nil, 0, e
	}
	var users []dto.UsersDto
	for _, user := range userRoleList {
		userDto := dto.UsersDto{}
		copier.Copy(&userDto, user)
		userDto.RoleIds = user.GetRoleIds()
		users = append(users, userDto)
	}
	return users, i, err
}

// 更新密码
func (us *UserService) ChangePwd(username string, newPasswd string) error {
	hashNewPasswd := util.GenPasswd(newPasswd)
	err := us.userDao.ChangePwd(username, hashNewPasswd)
	if err == nil {
		user := us.userDao.GetUserByUsername(username)
		if user != nil {
			cache.SetUser(user)
		}
	}
	return err
}

// 创建用户
func (us *UserService) CreateUser(user *model.User) error {
	user.Password = util.GenPasswd(user.Password)
	return database.Create(user)
}
func (us *UserService) GetUserByUsername(username string) *model.User {
	return us.userDao.GetUserByUsername(username)
}

// 更新用户
func (us *UserService) UpdateUser(user *model.User) error {
	err := us.userDao.UpdateUser(user)
	if err == nil {
		cache.SetUser(user)
	}
	return err
}

// 批量删除
func (us *UserService) BatchDeleteUserByIds(ids []uint64) error {
	users, err := us.userDao.GetUsersWithRoles(ids)
	if err != nil {
		return err
	}

	err = us.userDao.BatchDeleteUserByIds(ids)
	if err == nil {
		for _, user := range users {
			cache.UserCache.Delete(fmt.Sprintf("%v", user.ID))
		}
	}
	return err
}

// 根据用户ID获取用户角色排序最小值
func (us *UserService) GetUserMinRoleSortsByIds(ids []uint64) ([]int, error) {
	userList, err := us.userDao.GetUsersWithRoles(ids)
	if err != nil {
		return nil, err
	}

	roleMinSortList := make([]int, len(userList))
	for i, user := range userList {
		if len(user.Roles) == 0 {
			roleMinSortList[i] = 999
			continue
		}
		roleSortList := make([]int, len(user.Roles))
		for j, role := range user.Roles {
			roleSortList[j] = int(role.Sort)
		}
		roleMinSortList[i] = slices.Min(roleSortList)
	}
	return roleMinSortList, nil
}

// 设置用户信息缓存
//func (us *UserService) SetUserInfoCache(username string, user model.User) {
//	cache.Cache.Set(username, user, 0)
//}

// 根据角色ID更新拥有该角色的用户信息缓存
func (us *UserService) UpdateUserInfoCacheByRoleId(roleId uint64) error {
	roleDao := dao.RoleDao{}
	role, err := roleDao.GetRoleWithUsers(roleId)
	if err != nil {
		return errors.New("根据角色ID角色信息失败")
	}
	for _, user := range role.Users {
		found := cache.GetUser(user.ID)
		if found != nil {
			cache.SetUser(user)
			//cache.Cache.Set(user.UserName, *user, 0)
		}
	}
	return nil
}

// 清理所有用户信息缓存
//func (us *UserService) ClearUserInfoCache() {
//	cache.Cache.Flush()
//}
