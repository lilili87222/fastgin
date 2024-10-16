package dao

import (
	"codegen/output/model"
	"gorm.io/gorm"
)

type DictionaryDAO struct {
	db *gorm.DB
}

func NewDictionaryDAO(db *gorm.DB) *DictionaryDAO {
	return &DictionaryDAO{db: db}
}

func (dao *DictionaryDAO) Create(entity *model.Dictionary) error {
	return dao.db.Create(entity).Error
}

func (dao *DictionaryDAO) GetByID(id uint) (*model.Dictionary, error) {
	var entity model.Dictionary
	err := dao.db.First(&entity, id).Error
	return &entity, err
}

func (dao *DictionaryDAO) Update(entity *model.Dictionary) error {
	return dao.db.Save(entity).Error
}

func (dao *DictionaryDAO) Delete(id uint) error {
	return dao.db.Delete(&model.Dictionary{}, id).Error
}
