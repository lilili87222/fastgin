package service

import (
	"fastgin/common/httpz"
	"fastgin/database"
	"fastgin/modules/app/dao"
	"fastgin/modules/app/model"
)

type DictionaryService struct {
	dao *dao.DictionaryDAO
}

func NewDictionaryService() *DictionaryService {
	return &DictionaryService{dao: dao.NewDictionaryDAO()}
}

func (s *DictionaryService) Create(entity *model.Dictionary) error {
	return database.Create(entity)
}

func (s *DictionaryService) GetByID(id uint) (model.Dictionary, error) {
	return database.GetById[model.Dictionary](id)
}

func (s *DictionaryService) Update(entity *model.Dictionary) error {
	return database.Update(entity)
}

func (s *DictionaryService) Delete(id uint) error {
	return database.Delete[model.Dictionary](id)
}

func (s *DictionaryService) Search(req *httpz.SearchRequest) ([]model.Dictionary, int64, error) {
	return database.SearchTable[model.Dictionary](req)
}
func (s *DictionaryService) DeleteBatch(ids []uint) error {
	return database.DeleteByIds[model.Dictionary](ids)
}
