package dao

import (
	"fastgin/database"
	"fastgin/modules/sys/model"
)

type UserDao struct{}

func NewUserDao() *UserDao {
	return &UserDao{}
}

func (ur *UserDao) GetUserByUsername(username string) *model.User {
	var user model.User
	err := database.DB.Where("user_name = ?", username).Preload("Roles").First(&user).Error
	if err != nil {
		return nil
	}
	return &user
}

func (ur *UserDao) GetUserWithRoles(id uint64) (*model.User, error) {
	return database.GetByIdPreload[*model.User](id, "Roles")
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

func (ur *UserDao) BatchDeleteUserByIds(ids []uint64) error {
	users, e := ur.GetUsersWithRoles(ids)
	if e != nil {
		return e
	}
	return database.DB.Select("Roles").Unscoped().Delete(&users).Error
}

func (ur *UserDao) GetUsersWithRoles(ids []uint64) ([]*model.User, error) {
	var users []*model.User
	err := database.DB.Where("id IN (?)", ids).Preload("Roles").Find(&users).Error
	return users, err
}
