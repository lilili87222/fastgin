package service

import (
	"fastgin/common/httpz"
	"fastgin/database"
	"fastgin/modules/sys/dao"
	"fastgin/modules/sys/model"
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

func (s *DictionaryService) GetByID(id uint64) (model.Dictionary, error) {
	return database.GetById[model.Dictionary](id)
}

func (s *DictionaryService) Update(entity *model.Dictionary) error {
	return database.Update(entity)
}

func (s *DictionaryService) Delete(id uint64) error {
	return database.Delete[model.Dictionary](id)
}

func (s *DictionaryService) Search(req *httpz.SearchRequest) ([]model.Dictionary, int64, error) {
	return database.SearchTable[model.Dictionary](req)
}
func (s *DictionaryService) DeleteBatch(ids []uint64) error {
	return database.DeleteByIds[model.Dictionary](ids)
}
