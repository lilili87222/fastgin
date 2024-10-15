package dao

import (
	"fastgin/database"
	"fastgin/sys/dto"
	"fastgin/sys/model"
	"fmt"
	"strings"
)

type ApiDao struct {
}

func (a *ApiDao) GetApis(req *dto.ApiListRequest) ([]*model.Api, int64, error) {
	var list []*model.Api
	db := database.DB.Model(&model.Api{}).Order("created_at DESC")

	method := strings.TrimSpace(req.Method)
	if method != "" {
		db = db.Where("method LIKE ?", fmt.Sprintf("%%%s%%", method))
	}
	path := strings.TrimSpace(req.Path)
	if path != "" {
		db = db.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	category := strings.TrimSpace(req.Category)
	if category != "" {
		db = db.Where("category LIKE ?", fmt.Sprintf("%%%s%%", category))
	}
	creator := strings.TrimSpace(req.Creator)
	if creator != "" {
		db = db.Where("creator LIKE ?", fmt.Sprintf("%%%s%%", creator))
	}

	var total int64
	err := db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := int(req.PageNum)
	pageSize := int(req.PageSize)
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}

func (a *ApiDao) GetApisById(apiIds []uint) ([]*model.Api, error) {
	var apis []*model.Api
	err := database.DB.Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}

func (a *ApiDao) GetApiTree() ([]*model.Api, error) {
	var apiList []*model.Api
	err := database.DB.Order("category").Order("created_at").Find(&apiList).Error
	return apiList, err
}

func (a *ApiDao) CreateApi(api *model.Api) error {
	return database.DB.Create(api).Error
}

func (a *ApiDao) UpdateApiById(apiId uint, api *model.Api) error {
	return database.DB.Model(api).Where("id = ?", apiId).Updates(api).Error
}

func (a *ApiDao) BatchDeleteApiByIds(apiIds []uint) error {
	return database.DB.Where("id IN (?)", apiIds).Unscoped().Delete(&model.Api{}).Error
}

func (a *ApiDao) GetApiDescByPath(path string, method string) (string, error) {
	var api model.Api
	err := database.DB.Where("path = ?", path).Where("method = ?", method).First(&api).Error
	return api.Desc, err
}
