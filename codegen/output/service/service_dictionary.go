package service

import (
	"codegen/output/dao"
	"codegen/output/model"
)

type DictionaryService struct {
	dao *dao.DictionaryDAO
}

func NewDictionaryService(dao *dao.DictionaryDAO) *DictionaryService {
	return &DictionaryService{dao: dao}
}

func (s *DictionaryService) Create(entity *model.Dictionary) error {
	return s.dao.Create(entity)
}

func (s *DictionaryService) GetByID(id uint) (*model.Dictionary, error) {
	return s.dao.GetByID(id)
}

func (s *DictionaryService) Update(entity *model.Dictionary) error {
	return s.dao.Update(entity)
}

func (s *DictionaryService) Delete(id uint) error {
	return s.dao.Delete(id)
}
