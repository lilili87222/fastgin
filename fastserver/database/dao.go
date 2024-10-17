package database

import (
	"fastgin/sys/dto"
	"fmt"
)

func Create(itemPoint any) error {
	return DB.Create(itemPoint).Error
}
func Delete[T any](id uint) error {
	return DB.Delete(new(T), id).Error
}
func Update(api any) error {
	return DB.Model(api).Save(api).Error
}
func GetById[T any](id uint) (T, error) {
	var item T
	err := DB.First(&item, id).Error
	return item, err
}

func DeleteByIds[T any](ids []uint) error {
	return DB.Where("id IN (?)", ids).Unscoped().Delete(new(T)).Error
}

func GetByIds[T any](ids []uint) ([]T, error) {
	var apis []T
	err := DB.Where("id IN (?)", ids).Find(&apis).Error
	return apis, err
}

func GetByIdPreload[T any](id uint, preloads ...string) (T, error) {
	var item T
	db := DB
	for _, preload := range preloads {
		db = db.Preload(preload)
	}
	err := db.First(&item, id).Error
	return item, err
}

func ListAll[T any](orders ...string) ([]T, error) {
	var items []T
	db := DB
	for _, s := range orders {
		db = db.Order(s)
	}
	err := db.Find(&items).Error
	return items, err
}

func SearchTable[T any](req *dto.SearchRequest) ([]T, int64, error) {
	var list []T
	db := DB.Model(new(T)).Order("created_at DESC")
	for key, value := range req.KeyValues {
		if key == "status" {
			continue
		}
		db = db.Where(key+" LIKE ?", fmt.Sprintf("%%%s%%", value))
	}
	if req.KeyValues["status"] != "" {
		db = db.Where("status = ?", req.KeyValues["status"])
	}
	// 分页
	var total int64
	var err = db.Count(&total).Error
	if err != nil {
		return list, total, err
	}
	pageNum := req.PageNum
	pageSize := req.PageSize
	if pageNum > 0 && pageSize > 0 {
		err = db.Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&list).Error
	} else {
		err = db.Find(&list).Error
	}
	return list, total, err
}
