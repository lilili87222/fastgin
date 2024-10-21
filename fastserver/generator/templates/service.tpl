package service

import (
 "fastgin/common/httpz"
 "fastgin/database"
 "{{.Module}}/dao"
 "{{.Module}}/model"
)

type {{.ModelName}}Service struct {
 dao *dao.{{.ModelName}}DAO
}

func New{{.ModelName}}Service() *{{.ModelName}}Service {
 return &{{.ModelName}}Service{dao: dao.New{{.ModelName}}DAO()}
}

func (s *{{.ModelName}}Service) Create(entity *model.{{.ModelName}}) error {
 return database.Create(entity)
}

func (s *{{.ModelName}}Service) GetByID(id uint64) (model.{{.ModelName}}, error) {
 return database.GetById[model.{{.ModelName}}](id)
}

func (s *{{.ModelName}}Service) Update(entity *model.{{.ModelName}}) error {
 return database.Update(entity)
}

func (s *{{.ModelName}}Service) Delete(id uint64) error {
 return database.Delete[model.{{.ModelName}}](id)
}

func (s *{{.ModelName}}Service) Search(req *httpz.SearchRequest) ([]model.{{.ModelName}}, int64, error) {
 return database.SearchTable[model.{{.ModelName}}](req)
}
func (s *{{.ModelName}}Service) DeleteBatch(ids []uint64) error {
 return database.DeleteByIds[model.{{.ModelName}}](ids)
}
