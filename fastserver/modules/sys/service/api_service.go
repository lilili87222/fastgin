package service

import (
	"errors"
	"fastgin/config"
	"fastgin/database"
	"fastgin/modules/sys/dao"
	"fastgin/modules/sys/dto"
	"fastgin/modules/sys/model"
	"github.com/thoas/go-funk"
)

type ApiService struct {
	apiDao *dao.ApiDao
}

func NewApiService() *ApiService {
	return &ApiService{apiDao: &dao.ApiDao{}}
}

//func (s *ApiService) List(req *dto.SearchRequest) ([]model.Api, int64, error) {
//	return database.SearchTable[model.Api](req)
//}

func (s *ApiService) GetApiTree() ([]*dto.ApiTreeDto, error) {

	//apiList, err := s.apiDao.GetApiTree()
	apiList, err := database.ListAll[model.Api]("category", "created_at")
	if err != nil {
		return nil, err
	}

	var categoryList []string
	for _, api := range apiList {
		categoryList = append(categoryList, api.Category)
	}
	categoryUniq := funk.UniqString(categoryList)

	apiTree := make([]*dto.ApiTreeDto, len(categoryUniq))
	for i, category := range categoryUniq {
		apiTree[i] = &dto.ApiTreeDto{
			Id:       -i,
			Desc:     category,
			Category: category,
			Children: nil,
		}
		for _, api := range apiList {
			if category == api.Category {
				apiTree[i].Children = append(apiTree[i].Children, &api)
			}
		}
	}
	return apiTree, nil
}

func (s *ApiService) CreateApi(api *model.Api) error {
	return database.Create(api)
	//return s.apiDao.Create(api)
}

func (s *ApiService) UpdateApiById(api *model.Api) error {
	//oldApi, err := s.apiDao.GetApisById([]uint{apiId})
	oldApi, err := database.GetById[model.Api](api.Id)
	if err != nil {
		return errors.New("根据接口ID获取接口信息失败")
	}

	//err = s.apiDao.Update(apiId, api)
	err = database.Update(api)
	if err != nil {
		return err
	}

	if oldApi.Path != api.Path || oldApi.Method != api.Method {
		policies, err := config.CasbinEnforcer.GetFilteredPolicy(1, oldApi.Path, oldApi.Method)
		if err != nil {
			return err
		}
		if len(policies) > 0 {
			isRemoved, _ := config.CasbinEnforcer.RemovePolicies(policies)
			if !isRemoved {
				return errors.New("更新权限接口失败")
			}
			for _, policy := range policies {
				policy[1] = api.Path
				policy[2] = api.Method
			}
			isAdded, _ := config.CasbinEnforcer.AddPolicies(policies)
			if !isAdded {
				return errors.New("更新权限接口失败")
			}
			err = config.CasbinEnforcer.LoadPolicy()
			if err != nil {
				return errors.New("更新权限接口成功，权限接口策略加载失败")
			}
		}
	}
	return nil
}

func (s *ApiService) BatchDeleteApiByIds(apiIds []uint) error {
	//apis, err := s.apiDao.GetApisById(apiIds)
	apis, err := database.GetByIds[model.Api](apiIds)
	if err != nil {
		return errors.New("根据接口ID获取接口列表失败")
	}
	if len(apis) == 0 {
		return errors.New("根据接口ID未获取到接口列表")
	}

	err = database.DeleteByIds[model.Api](apiIds)
	//err = s.apiDao.BatchDelete(apiIds)
	if err == nil {
		for _, api := range apis {
			policies, err := config.CasbinEnforcer.GetFilteredPolicy(1, api.Path, api.Method)
			if err != nil {
				return err
			}
			if len(policies) > 0 {
				isRemoved, _ := config.CasbinEnforcer.RemovePolicies(policies)
				if !isRemoved {
					return errors.New("删除权限接口失败")
				}
			}
		}
		err = config.CasbinEnforcer.LoadPolicy()
		if err != nil {
			return errors.New("删除权限接口成功，权限接口策略加载失败")
		}
	}
	return err
}

func (s *ApiService) GetApiDescByPath(path string, method string) (model.Api, error) {
	return s.apiDao.GetApiDescByPath(path, method)
}

func (s *ApiService) GetApisById(apiIds []uint) ([]model.Api, error) {
	return database.GetByIds[model.Api](apiIds)
}
