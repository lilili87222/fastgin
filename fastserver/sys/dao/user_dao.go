package dao

import (
	"fastgin/config"
	"fastgin/sys/dto"
	"fastgin/sys/model"
	"fmt"
	"strings"
)

type UserDao struct{}

func NewUserDao() UserDao {
	return UserDao{}
}

func (ur UserDao) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := config.DB.Where("username = ?", username).Preload("Roles").First(&user).Error
	return user, err
}

func (ur UserDao) GetUserById(id uint) (model.User, error) {
	var user model.User
	err := config.DB.Where("id = ?", id).Preload("Roles").First(&user).Error
	return user, err
}

func (ur UserDao) GetUsers(req *dto.UserListRequest) ([]*model.User, int64, error) {
	var list []*model.User
	db := config.DB.Model(&model.User{}).Order("created_at DESC")

	if username := strings.TrimSpace(req.Username); username != "" {
		db = db.Where("username LIKE ?", fmt.Sprintf("%%%s%%", username))
	}
	if nickname := strings.TrimSpace(req.Nickname); nickname != "" {
		db = db.Where("nickname LIKE ?", fmt.Sprintf("%%%s%%", nickname))
	}
	if mobile := strings.TrimSpace(req.Mobile); mobile != "" {
		db = db.Where("mobile LIKE ?", fmt.Sprintf("%%%s%%", mobile))
	}
	if status := req.Status; status != 0 {
		db = db.Where("status = ?", status)
	}

	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}

	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Preload("Roles").Find(&list).Error
	} else {
		err = db.Preload("Roles").Find(&list).Error
	}
	return list, total, err
}

func (ur UserDao) ChangePwd(username string, hashNewPasswd string) error {
	return config.DB.Model(&model.User{}).Where("username = ?", username).Update("password", hashNewPasswd).Error
}

func (ur UserDao) CreateUser(user *model.User) error {
	return config.DB.Create(user).Error
}

func (ur UserDao) UpdateUser(user *model.User) error {
	err := config.DB.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	return config.DB.Model(user).Association("Roles").Replace(user.Roles)
}

func (ur UserDao) BatchDeleteUserByIds(ids []uint) error {
	users, e := ur.GetUsersByIds(ids)
	if e != nil {
		return e
	}
	return config.DB.Select("Roles").Unscoped().Delete(&users).Error
}

func (ur UserDao) GetUsersByIds(ids []uint) ([]model.User, error) {
	var users []model.User
	err := config.DB.Where("id IN (?)", ids).Preload("Roles").Find(&users).Error
	return users, err
}

func (ur UserDao) GetRoleById(roleId uint) (model.Role, error) {
	var role model.Role
	err := config.DB.Where("id = ?", roleId).Preload("Users").First(&role).Error
	return role, err
}
