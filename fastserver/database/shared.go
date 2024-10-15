package database

import (
	"fastgin/sys/dto"
	"fmt"
)

func SearchTable[T any](req *dto.SearchRequest) ([]T, int64, error) {
	var list []T
	db := DB.Model(new(T)).Order("start_time DESC")
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
	err := db.Count(&total).Error
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
func DeleteByIds[T any](ids []uint) error {
	return DB.Where("id IN (?)", ids).Unscoped().Delete(new(T)).Error
}
