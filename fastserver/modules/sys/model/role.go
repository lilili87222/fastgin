package model

//	type Role struct {
//		gorm.Model
//		Name    string  `gorm:"type:varchar(20);not null;unique" json:"name"`
//		Keyword string  `gorm:"type:varchar(20);not null;unique" json:"keyword"`
//		Desc    *string `gorm:"type:varchar(100);" json:"desc"`
//		Status  uint    `gorm:"type:tinyint(1);default:1;comment:'1正常, 2禁用'" json:"status"`
//		Sort    uint    `gorm:"type:int(3);default:999;comment:'角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)'" json:"sort"`
//		Creator string  `gorm:"type:varchar(20);" json:"creator"`
//		Users   []*User `gorm:"many2many:sys_user_role" json:"users"`
//		Menus   []*Menu `gorm:"many2many:sys_role_menu;" json:"menus"` // 角色菜单多对多关系
//	}
type Role struct {
	Model
	Name    string  `gorm:"type:varchar(20);not null;unique" `
	Keyword string  `gorm:"type:varchar(20);not null;unique" `
	Desc    *string `gorm:"type:varchar(100);" `
	Status  uint    `gorm:"type:tinyint(1);default:1;comment:'1正常, 2禁用'" `
	Sort    uint    `gorm:"type:int(3);default:999;comment:'角色排序(排序越大权限越低, 不能查看比自己序号小的角色, 不能编辑同序号用户权限, 排序为1表示超级管理员)'" `
	Creator string  `gorm:"type:varchar(20);" `
	Users   []*User `gorm:"many2many:sys_user_role" `
	Menus   []*Menu `gorm:"many2many:sys_role_menu;" ` // 角色菜单多对多关系
}

func (*Role) TableName() string {
	return "sys_role"
}
