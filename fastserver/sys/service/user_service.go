package service

import (
	"errors"
	"fastgin/sys/dao"
	"fastgin/sys/dto"
	"fastgin/sys/model"
	"fastgin/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/patrickmn/go-cache"
	"slices"
	"time"
)

type UserService struct {
	userDao *dao.UserDao
}

var userInfoCache = cache.New(24*time.Hour, 48*time.Hour)

func NewUserService() *UserService {
	return &UserService{userDao: dao.NewUserDao()}
}

// 登录
func (us *UserService) Login(user *model.User) (*model.User, error) {
	firstUser, err := us.userDao.GetUserByUsername(user.UserName)
	if err != nil {
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

	err = util.ComparePasswd(firstUser.Password, user.Password)
	if err != nil {
		return &firstUser, errors.New("密码错误")
	}
	return &firstUser, nil
}

// 获取当前登录用户信息
func (us *UserService) GetCurrentUser(c *gin.Context) (model.User, error) {
	ctxUser, exist := c.Get("user")
	if !exist {
		return model.User{}, errors.New("用户未登录")
	}
	u, _ := ctxUser.(model.User)

	cacheUser, found := userInfoCache.Get(u.UserName)
	if found {
		return cacheUser.(model.User), nil
	}

	user, err := us.userDao.GetUserById(u.Id)
	if err != nil {
		userInfoCache.Delete(u.UserName)
	} else {
		userInfoCache.Set(u.UserName, user, cache.DefaultExpiration)
	}
	return user, err
}

// 获取当前用户角色排序最小值（最高等级角色）以及当前用户信息
func (us *UserService) GetCurrentUserMinRoleSort(c *gin.Context) (uint, model.User, error) {
	ctxUser, err := us.GetCurrentUser(c)
	if err != nil {
		return 999, ctxUser, err
	}

	currentRoleSorts := make([]int, len(ctxUser.Roles))
	for i, role := range ctxUser.Roles {
		currentRoleSorts[i] = int(role.Sort)
	}

	currentRoleSortMin := slices.Min(currentRoleSorts)
	return uint(currentRoleSortMin), ctxUser, nil
}

// 获取单个用户
func (us *UserService) GetUserById(id uint) (model.User, error) {
	return us.userDao.GetUserById(id)
}

// 获取用户列表
func (us *UserService) GetUsers(req *dto.UserListRequest) ([]*model.User, int64, error) {
	return us.userDao.GetUsers(req)
}
func (us *UserService) GetUsersWithRoleIds(req *dto.UserListRequest) ([]dto.UsersDto, int64, error) {
	userList, i, err := us.userDao.GetUsers(req)
	var users []dto.UsersDto
	for _, user := range userList {
		userDto := dto.UsersDto{}
		copier.Copy(&userDto, user)
		userDto.RoleIds = user.GetRoleIds()
		users = append(users, userDto)
	}
	return users, i, err
}

// 更新密码
func (us *UserService) ChangePwd(username string, hashNewPasswd string) error {
	err := us.userDao.ChangePwd(username, hashNewPasswd)
	if err == nil {
		cacheUser, found := userInfoCache.Get(username)
		if found {
			user := cacheUser.(model.User)
			user.Password = hashNewPasswd
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		} else {
			user, _ := us.userDao.GetUserByUsername(username)
			userInfoCache.Set(username, user, cache.DefaultExpiration)
		}
	}
	return err
}

// 创建用户
func (us *UserService) CreateUser(user *model.User) error {
	return us.userDao.CreateUser(user)
}

// 更新用户
func (us *UserService) UpdateUser(user *model.User) error {
	err := us.userDao.UpdateUser(user)
	if err == nil {
		userInfoCache.Set(user.UserName, *user, cache.DefaultExpiration)
	}
	return err
}

// 批量删除
func (us *UserService) BatchDeleteUserByIds(ids []uint) error {
	users, err := us.userDao.GetUsersByIds(ids)
	if err != nil {
		return err
	}

	err = us.userDao.BatchDeleteUserByIds(ids)
	if err == nil {
		for _, user := range users {
			userInfoCache.Delete(user.UserName)
		}
	}
	return err
}

// 根据用户ID获取用户角色排序最小值
func (us *UserService) GetUserMinRoleSortsByIds(ids []uint) ([]int, error) {
	userList, err := us.userDao.GetUsersByIds(ids)
	if err != nil {
		return nil, err
	}

	roleMinSortList := make([]int, len(userList))
	for i, user := range userList {
		roleSortList := make([]int, len(user.Roles))
		for j, role := range user.Roles {
			roleSortList[j] = int(role.Sort)
		}
		roleMinSortList[i] = slices.Min(roleSortList)
	}
	return roleMinSortList, nil
}

// 设置用户信息缓存
func (us *UserService) SetUserInfoCache(username string, user model.User) {
	userInfoCache.Set(username, user, cache.DefaultExpiration)
}

// 根据角色ID更新拥有该角色的用户信息缓存
func (us *UserService) UpdateUserInfoCacheByRoleId(roleId uint) error {
	role, err := us.userDao.GetRoleById(roleId)
	if err != nil {
		return errors.New("根据角色ID角色信息失败")
	}

	for _, user := range role.Users {
		_, found := userInfoCache.Get(user.UserName)
		if found {
			userInfoCache.Set(user.UserName, *user, cache.DefaultExpiration)
		}
	}
	return nil
}

// 清理所有用户信息缓存
func (us *UserService) ClearUserInfoCache() {
	userInfoCache.Flush()
}
