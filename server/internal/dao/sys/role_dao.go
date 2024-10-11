package sys

import (
	"errors"
	"fmt"
	"go-web-mini/config"
	"go-web-mini/internal/bean"
	"go-web-mini/internal/model/sys"
	"strings"
)

type IRoleRepository interface {
	GetRoles(req *bean.RoleListRequest) ([]sys.Role, int64, error)       // 获取角色列表
	GetRolesByIds(roleIds []uint) ([]*sys.Role, error)                   // 根据角色ID获取角色
	CreateRole(role *sys.Role) error                                     // 创建角色
	UpdateRoleById(roleId uint, role *sys.Role) error                    // 更新角色
	GetRoleMenusById(roleId uint) ([]*sys.Menu, error)                   // 获取角色的权限菜单
	UpdateRoleMenus(role *sys.Role) error                                // 更新角色的权限菜单
	GetRoleApisByRoleKeyword(roleKeyword string) ([]*sys.Api, error)     // 根据角色关键字获取角色的权限接口
	UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error // 更新角色的权限接口（先全部删除再新增）
	BatchDeleteRoleByIds(roleIds []uint) error                           // 删除角色
}

type RoleRepository struct {
}

func NewRoleRepository() IRoleRepository {
	return RoleRepository{}
}

// 获取角色列表
func (r RoleRepository) GetRoles(req *bean.RoleListRequest) ([]sys.Role, int64, error) {
	var list []sys.Role
	db := config.DB.Model(&sys.Role{}).Order("created_at DESC")

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
func (r RoleRepository) GetRolesByIds(roleIds []uint) ([]*sys.Role, error) {
	var list []*sys.Role
	err := config.DB.Where("id IN (?)", roleIds).Find(&list).Error
	return list, err
}

// 创建角色
func (r RoleRepository) CreateRole(role *sys.Role) error {
	err := config.DB.Create(role).Error
	return err
}

// 更新角色
func (r RoleRepository) UpdateRoleById(roleId uint, role *sys.Role) error {
	err := config.DB.Model(&sys.Role{}).Where("id = ?", roleId).Updates(role).Error
	return err
}

// 获取角色的权限菜单
func (r RoleRepository) GetRoleMenusById(roleId uint) ([]*sys.Menu, error) {
	var role sys.Role
	err := config.DB.Where("id = ?", roleId).Preload("Menus").First(&role).Error
	return role.Menus, err
}

// 更新角色的权限菜单
func (r RoleRepository) UpdateRoleMenus(role *sys.Role) error {
	err := config.DB.Model(role).Association("Menus").Replace(role.Menus)
	return err
}

// 根据角色关键字获取角色的权限接口
func (r RoleRepository) GetRoleApisByRoleKeyword(roleKeyword string) ([]*sys.Api, error) {
	policies := config.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)

	// 获取所有接口
	var apis []*sys.Api
	err := config.DB.Find(&apis).Error
	if err != nil {
		return apis, errors.New("获取角色的权限接口失败")
	}

	accessApis := make([]*sys.Api, 0)

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
func (r RoleRepository) UpdateRoleApis(roleKeyword string, reqRolePolicies [][]string) error {
	// 先获取path中的角色ID对应角色已有的police(需要先删除的)
	err := config.CasbinEnforcer.LoadPolicy()
	if err != nil {
		return errors.New("角色的权限接口策略加载失败")
	}
	rmPolicies := config.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
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
func (r RoleRepository) BatchDeleteRoleByIds(roleIds []uint) error {
	var roles []*sys.Role
	err := config.DB.Where("id IN (?)", roleIds).Find(&roles).Error
	if err != nil {
		return err
	}
	err = config.DB.Select("Users", "Menus").Unscoped().Delete(&roles).Error
	// 删除成功就删除casbin policy
	if err == nil {
		for _, role := range roles {
			roleKeyword := role.Keyword
			rmPolicies := config.CasbinEnforcer.GetFilteredPolicy(0, roleKeyword)
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
