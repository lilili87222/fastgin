package dto

// 创建接口结构体
type CreateMenuRequest struct {
	Name       string `json:"name" form:"name" validate:"required,min=1,max=50"`
	Title      string `json:"title" form:"title" validate:"required,min=1,max=50"`
	Icon       string `json:"icon" form:"icon" validate:"min=0,max=50"`
	Path       string `json:"path" form:"path" validate:"required,min=1,max=100"`
	Redirect   string `json:"redirect" form:"redirect" validate:"min=0,max=100"`
	Component  string `json:"component" form:"component" validate:"required,min=1,max=100"`
	Sort       uint   `json:"sort" form:"sort" validate:"gte=1,lte=999"`
	Status     uint   `json:"status" form:"status" validate:"oneof=1 2"`
	Hidden     uint   `json:"hidden" form:"hidden" validate:"oneof=1 2"`
	NoCache    uint   `json:"no_cache" form:"no_cache" validate:"oneof=1 2"`
	AlwaysShow uint   `json:"always_show" form:"always_show" validate:"oneof=1 2"`
	Breadcrumb uint   `json:"breadcrumb" form:"breadcrumb" validate:"oneof=1 2"`
	ActiveMenu string `json:"active_menu" form:"active_menu" validate:"min=0,max=100"`
	ParentId   uint64 `json:"parent_id" form:"parent_id"`
}
