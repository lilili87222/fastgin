package dao

import (
	"fastgin/database"
	"fastgin/modules/sys/model"
)

type MenuDao struct {
}

// 批量删除菜单
func (m *MenuDao) BatchDeleteMenuByIds(menuIds []uint64) error {
	menus, err := database.GetByIds[model.Menu](menuIds)
	err = database.DB.Select("Roles").Unscoped().Delete(&menus).Error
	return err
}
