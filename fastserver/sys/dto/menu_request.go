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
	ParentId   uint   `json:"ParentId" form:"parentId"`
}

// 更新接口结构体
type UpdateMenuRequest struct {
	Name       string `json:"Name" form:"Name" validate:"required,min=1,max=50"`
	Title      string `json:"Title" form:"Title" validate:"required,min=1,max=50"`
	Icon       string `json:"Icon" form:"Icon" validate:"min=0,max=50"`
	Path       string `json:"Path" form:"Path" validate:"required,min=1,max=100"`
	Redirect   string `json:"Redirect" form:"Redirect" validate:"min=0,max=100"`
	Component  string `json:"Component" form:"Component" validate:"min=0,max=100"`
	Sort       uint   `json:"Sort" form:"Sort" validate:"gte=1,lte=999"`
	Status     uint   `json:"Status" form:"Status" validate:"oneof=1 2"`
	Hidden     uint   `json:"Hidden" form:"Hidden" validate:"oneof=1 2"`
	NoCache    uint   `json:"NoCache" form:"NoCache" validate:"oneof=1 2"`
	AlwaysShow uint   `json:"AlwaysShow" form:"AlwaysShow" validate:"oneof=1 2"`
	Breadcrumb uint   `json:"Breadcrumb" form:"Breadcrumb" validate:"oneof=1 2"`
	ActiveMenu string `json:"ActiveMenu" form:"ActiveMenu" validate:"min=0,max=100"`
	ParentId   uint   `json:"ParentId" form:"ParentId"`
}

// 删除接口结构体
type DeleteMenuRequest struct {
	MenuIds []uint `json:"MenuIds" form:"MenuIds"`
}
