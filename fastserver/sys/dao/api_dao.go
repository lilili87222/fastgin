package dao

import (
	"fastgin/database"
	"fastgin/sys/model"
)

type ApiDao struct {
}

func (a *ApiDao) GetApiDescByPath(path string, method string) (string, error) {
	var api model.Api
	err := database.DB.Where("path = ?", path).Where("method = ?", method).First(&api).Error
	return api.Desc, err
}
