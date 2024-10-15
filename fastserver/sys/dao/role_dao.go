package dao

import (
	"fastgin/config"
	"fastgin/sys/model"
	"fastgin/sys/model/request"
	"fmt"
	"strings"
)

type RoleDao struct {
}

// 获取角色列表
func (r RoleDao) GetRoles(req *request.RoleListRequest) ([]model.Role, int64, error) {
	var list []model.Role
	db := config.DB.Model(&model.Role{}).Order("created_at DESC")

	name := strings.TrimSpace(req.Name)
	if name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%%%s%%", name))
	}
	keyword := strings.TrimSpace(req.Keyword)
	if keyword != "" {
		db = db.Where("keyword LIKE ?", fmt.Sprintf("%%%s%%", keyword))
	}
	status := req.Status
	if status != 0 {
		db = db.Where("status = ?", status)
	}
	// 当pageNum > 0 且 pageSize > 0 才分页
	//记录总条数
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

// 根据角色ID获取角色
func (r RoleDao) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	var list []*model.Role
	err := config.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// 创建角色
func (r RoleDao) CreateRole(role *model.Role) error {
	err := config.DB.Create(role).Error
	return err
}

// 更新角色
func (r RoleDao) UpdateRoleById(roleId uint, role *model.Role) error {
	err := config.DB.Model(&model.Role{}).Where("id = ?", roleId).Updates(role).Error
	return err
}

// 获取角色的权限菜单
func (r RoleDao) GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	var role model.Role
	err := config.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// 更新角色的权限菜单
func (r RoleDao) UpdateRoleMenus(role *model.Role) error {
	err := config.DB.Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// 根据角色关键字获取角色的权限接口
//func (r RoleDao) GetRoleApisByRoleKeyword(roleKeyword string) ([]*model.Api, error) {
//	policies, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
//	if err2 != nil {
//		return nil, errors.New("获取角色的权限接口失败")
//	}
//
//	// 获取所有接口
//	var apis []*model.Api
//	err := config.DB.Find(&apis).Error
//	if err != nil {
//		return apis, errors.New("获取角色的权限接口失败")
//	}
//
//	accessApis := make([]*model.Api, 0)
//
//	for _, policy := range policies {
//		path := policy[1]
//		method := policy[2]
//		for _, api := range apis {
//			if path == api.Path && method == api.Method {
//				accessApis = append(accessApis, api)
//				break
//			}
//		}
//	}
//
//	return accessApis, err
//
//}

// 更新角色的权限接口（先全部删除再新增）
//func (r RoleDao) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
//	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
//	err := config.CasbinEnforcer.LoadPolicy()
//	if err != nil {
//		return errors.New("角色的权限接口策略加载失败")
//	}
//	rmPolicies, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
//	if err2 != nil {
//		return errors.New("获取角色的权限接口失败")
//	}
//	if len(rmPolicies) > 0 {
//		isRemoved, _ := config.CasbinEnforcer.RemovePolicies(rmPolicies)
//		if !isRemoved {
//			return errors.New("更新角色的权限接口失败")
//		}
//	}
//	isAdded, _ := config.CasbinEnforcer.AddPolicies(reqRolePolicies)
//	if !isAdded {
//		return errors.New("更新角色的权限接口失败")
//	}
//	err = config.CasbinEnforcer.LoadPolicy()
//	if err != nil {
//		return errors.New("更新角色的权限接口成功，角色的权限接口策略加载失败")
//	} else {
//		return err
//	}
//}

// 删除角色
func (r RoleDao) BatchDeleteRoleByIds(roleIds []uint) error {
	return nil
}
