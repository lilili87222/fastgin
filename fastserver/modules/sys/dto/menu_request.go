package dto

// 创建接口结构体
type CreateMenuRequest struct {
	Name       string `json:"Name" form:"name" validate:"required,min=1,max=50"`
	Title      string `json:"Title" form:"title" validate:"required,min=1,max=50"`
	Icon       string `json:"Icon" form:"icon" validate:"min=0,max=50"`
	Path       string `json:"Path" form:"path" validate:"required,min=1,max=100"`
	Redirect   string `json:"Redirect" form:"redirect" validate:"min=0,max=100"`
	Component  string `json:"Component" form:"component" validate:"required,min=1,max=100"`
	Sort       uint   `json:"Sort" form:"sort" validate:"gte=1,lte=999"`
	Status     uint   `json:"Status" form:"status" validate:"oneof=1 2"`
	Hidden     uint   `json:"Hidden" form:"hidden" validate:"oneof=1 2"`
	NoCache    uint   `json:"NoCache" form:"noCache" validate:"oneof=1 2"`
	AlwaysShow uint   `json:"AlwaysShow" form:"alwaysShow" validate:"oneof=1 2"`
	Breadcrumb uint   `json:"Breadcrumb" form:"breadcrumb" validate:"oneof=1 2"`
	ActiveMenu string `json:"ActiveMenu" form:"activeMenu" validate:"min=0,max=100"`
	ParentId   uint64 `json:"ParentId" form:"parentId"`
}
