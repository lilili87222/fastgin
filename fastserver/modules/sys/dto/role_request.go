package dto

// 新增角色结构体
type CreateRoleRequest struct {
	Name    string `json:"name" form:"name" validate:"required,min=1,max=20"`
	Keyword string `json:"keyword" form:"keyword" validate:"required,min=1,max=20"`
	Des     string `json:"des" form:"des" validate:"min=0,max=100"`
	Status  uint   `json:"status" form:"status" validate:"oneof=1 2"`
	Sort    int32  `json:"sort" form:"sort" validate:"gte=1,lte=999"`
}

// 获取用户角色结构体
type RoleListRequest struct {
	Name     string `json:"name" form:"name"`
	Keyword  string `json:"keyword" form:"keyword"`
	Status   uint   `json:"status" form:"status"`
	PageNum  uint   `json:"page_num" form:"page_num"`
	PageSize uint   `json:"page_size" form:"page_size"`
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
