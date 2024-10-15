package service

import (
	"errors"
	"fastgin/config"
	"fastgin/database"
	"fastgin/sys/dao"
	"fastgin/sys/dto"
	"fastgin/sys/model"
)

type RoleService struct {
	roleDao *dao.RoleDao
}

func NewRoleService() *RoleService {
	return &RoleService{roleDao: &dao.RoleDao{}}
}

// 获取角色列表
func (r *RoleService) GetRoles(req *dto.RoleListRequest) ([]model.Role, int64, error) {
	return r.roleDao.GetRoles(req)
}

// 根据角色ID获取角色
func (r *RoleService) GetRolesByIds(roleIds []uint) ([]*model.Role, error) {
	return r.roleDao.GetRolesByIds(roleIds)
}

// 创建角色
func (r *RoleService) CreateRole(role *model.Role) error {
	return r.roleDao.CreateRole(role)
}

// 更新角色
func (r *RoleService) UpdateRoleById(roleId uint, role *model.Role) error {
	return r.roleDao.UpdateRoleById(roleId, role)
}

// 获取角色的权限菜单
func (r *RoleService) GetRoleMenusById(roleId uint) ([]*model.Menu, error) {
	return r.roleDao.GetRoleMenusById(roleId)
}

// 更新角色的权限菜单
func (r *RoleService) UpdateRoleMenus(role *model.Role) error {
	return r.roleDao.UpdateRoleMenus(role)
}

// 根据角色关键字获取角色的权限接口
func (r *RoleService) GetRoleApisByRoleKeyword(roleKeyword string) ([]*model.Api, error) {
	policies, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if err2 != nil {
		return nil, errors.New("获取角色的权限接口失败")
	}

	// 获取所有接口
	var apis []*model.Api
	err := database.DB.Find(&apis).Error
	if err != nil {
		return apis, errors.New("获取角色的权限接口失败")
	}

	accessApis := make([]*model.Api, 0)

	for _, policy := range policies {
		path := policy[1]
		method := policy[2]
		for _, api := range apis {
			if path == api.Path && method == api.Method {
				accessApis = append(accessApis, api)
				break
			}
		}
	}

	return accessApis, err

}

// 更新角色的权限接口（先全部删除再新增）
func (r *RoleService) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := config.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("角色的权限接口策略加载失败")
	}
	rmPolicies, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
	if err2 != nil {
		return errors.New("获取角色的权限接口失败")
	}
	if len(rmPolicies) > 0 {
		isRemoved, _ := config.CasbinEnforcer.RemovePolicies(rmPolicies)
		if !isRemoved {
			return errors.New("更新角色的权限接口失败")
		}
	}
	isAdded, _ := config.CasbinEnforcer.AddPolicies(reqRolePolicies)
	if !isAdded {
		return errors.New("更新角色的权限接口失败")
	}
	err = config.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("更新角色的权限接口成功，角色的权限接口策略加载失败")
	} else {
		return err
	}
}

// 删除角色
func (r *RoleService) BatchDeleteRoleByIds(roleIds []uint) error {
	roles, err := r.roleDao.BatchDeleteRoleByIds(roleIds)
	// 删除成功就删除casbin policy
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies, err2 := config.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
			if err2 != nil {
				return errors.New("删除角色成功, 获取角色关联权限接口失败")
			}
			if len(rmPolicies) > 0 {
				isRemoved, _ := config.CasbinEnforcer.RemovePolicies(rmPolicies)
				if !isRemoved {
					return errors.New("删除角色成功, 删除角色关联权限接口失败")
				}
			}
		}
	}
	return err
}
