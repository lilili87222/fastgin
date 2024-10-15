package dto

// 新增角色结构体
type CreateRoleRequest struct {
	Name    string `json:"Name" form:"Name" validate:"required,min=1,max=20"`
	Keyword string `json:"Keyword" form:"Keyword" validate:"required,min=1,max=20"`
	Desc    string `json:"Desc" form:"Desc" validate:"min=0,max=100"`
	Status  uint   `json:"Status" form:"Status" validate:"oneof=1 2"`
	Sort    uint   `json:"Sort" form:"Sort" validate:"gte=1,lte=999"`
}

// 获取用户角色结构体
type RoleListRequest struct {
	Name     string `json:"Name" form:"Name"`
	Keyword  string `json:"Keyword" form:"Keyword"`
	Status   uint   `json:"Status" form:"Status"`
	PageNum  uint   `json:"PageNum" form:"PageNum"`
	PageSize uint   `json:"PageSize" form:"PageSize"`
}

// 批量删除角色结构体
//type DeleteRoleRequest struct {
//	RoleIds []uint `json:"RoleIds" form:"RoleIds"`
//}

// 更新角色的权限菜单
//type UpdateRoleMenusRequest struct {
//	MenuIds []uint `json:"MenuIds" form:"MenuIds"`
//}

// 更新角色的权限接口
//type UpdateRoleApisRequest struct {
//	ApiIds []uint `json:"ApiIds" form:"ApiIds"`
//}
