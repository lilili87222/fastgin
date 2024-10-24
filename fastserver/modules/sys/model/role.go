// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameRole = "sys_role"

// Role mapped from table <sys_role>
type Role struct {
	ID        uint64         `gorm:"column:id;type:bigint(20) unsigned;primaryKey;autoIncrement:true" json:"id"`
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime(3);default:CURRENT_TIMESTAMP(3)" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:datetime(3);default:CURRENT_TIMESTAMP(3)" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime(3)" json:"deleted_at"`
	Name      string         `gorm:"column:name;type:varchar(20);not null" json:"name"`
	Keyword   string         `gorm:"column:keyword;type:varchar(20);not null" json:"keyword"`
	Des       string         `gorm:"column:des;type:varchar(100)" json:"des"`
	Status    uint           `gorm:"column:status;type:tinyint(1);default:1;comment:'1正常, 2禁用'" json:"status"`                                          // '1正常, 2禁用'
	Sort      int32          `gorm:"column:sort;type:int(3);default:999;comment:'角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)'" json:"sort"` // '角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)'
	Creator   string         `gorm:"column:creator;type:varchar(20)" json:"creator"`
	Users     []*User        `gorm:"many2many:sys_user_role" json:"users"`
	Menus     []*Menu        `gorm:"many2many:sys_role_menu;" json:"menus" ` // 角色菜单多对多关系
}

// TableName Role's table name
func (*Role) TableName() string {
	return TableNameRole
}
