package config

import (
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// 全局CasbinEnforcer
var CasbinEnforcer *casbin.Enforcer

func InitCasbinEnforcer(db *gorm.DB) (*casbin.Enforcer, error) {
	a, err := gormadapter.NewAdapterByDBUseTableName(db, "sys", "casbin_rule")
	if err != nil {
		return nil, err
	}
	CasbinEnforcer, err = casbin.NewEnforcer(Instance.Casbin.ModelPath, a)
	if err != nil {
		return nil, err
	}

	err = CasbinEnforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}
	return CasbinEnforcer, nil
}
