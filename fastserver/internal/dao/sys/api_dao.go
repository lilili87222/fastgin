package sys

import (
	"fastgin/config"
	"fastgin/internal/model/sys"
	"fastgin/internal/model/sys/request"
	"fmt"
	"strings"
)

type ApiDao struct {
}

func (a ApiDao) GetApis(req *request.ApiListRequest) ([]*sys.Api, int64, error) {
	var list []*sys.Api
	db := config.DB.Model(&sys.Api{}).Order("created_at DESC")

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

func (a ApiDao) GetApisById(apiIds []uint) ([]*sys.Api, error) {
	var apis []*sys.Api
	err := config.DB.Where("id IN (?)", apiIds).Find(&apis).Error
	return apis, err
}

func (a ApiDao) GetApiTree() ([]*sys.Api, error) {
	var apiList []*sys.Api
	err := config.DB.Order("category").Order("created_at").Find(&apiList).Error
	return apiList, err
}

func (a ApiDao) CreateApi(api *sys.Api) error {
	return config.DB.Create(api).Error
}

func (a ApiDao) UpdateApiById(apiId uint, api *sys.Api) error {
	return config.DB.Model(api).Where("id = ?", apiId).Updates(api).Error
}

func (a ApiDao) BatchDeleteApiByIds(apiIds []uint) error {
	return config.DB.Where("id IN (?)", apiIds).Unscoped().Delete(&sys.Api{}).Error
}

func (a ApiDao) GetApiDescByPath(path string, method string) (string, error) {
	var api sys.Api
	err := config.DB.Where("path = ?", path).Where("method = ?", method).First(&api).Error
	return api.Desc, err
}
