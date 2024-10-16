package dao

import (
	"fastgin/database"
	"fastgin/sys/model"
)

type UserDao struct{}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (ur *UserDao) GetUserByUsername(username string) (model.User, error) {
	var user model.User
	err := database.DB.Where("user_name = ?", username).Preload("Roles").First(&user).Error
	return user, err
}

func (ur *UserDao) GetUserWithRoles(id uint) (model.User, error) {
	return database.GetByIdPreload[model.User](id, "Roles")
}
func (ur *UserDao) ChangePwd(username string, hashNewPasswd string) error {
	return database.DB.Model(&model.User{}).Where("user_name = ?", username).Update("password", hashNewPasswd).Error
}

func (ur *UserDao) UpdateUser(user *model.User) error {
	err := database.DB.Model(user).Updates(user).Error
	if err != nil {
		return err
	}
	return database.DB.Model(user).Association("Roles").Replace(user.Roles)
}

func (ur *UserDao) BatchDeleteUserByIds(ids []uint) error {
	users, e := ur.GetUsersWithRoles(ids)
	if e != nil {
		return e
	}
	return database.DB.Select("Roles").Unscoped().Delete(&users).Error
}

func (ur *UserDao) GetUsersWithRoles(ids []uint) ([]model.User, error) {
	var users []model.User
	err := database.DB.Where("id IN (?)", ids).Preload("Roles").Find(&users).Error
	return users, err
}
