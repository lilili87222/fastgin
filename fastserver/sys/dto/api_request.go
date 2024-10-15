package dto

import "fastgin/sys/model"

// 获取接口列表结构体
type ApiListRequest struct {
	Method   string `json:"method" form:"method"`
	Path     string `json:"path" form:"path"`
	Category string `json:"category" form:"category"`
	Creator  string `json:"creator" form:"creator"`
	PageNum  uint   `json:"pageNum" form:"pageNum"`
	PageSize uint   `json:"pageSize" form:"pageSize"`
}

// 创建接口结构体
type CreateApiRequest struct {
	Method   string `json:"method" form:"method" validate:"required,min=1,max=20"`
	Path     string `json:"path" form:"path" validate:"required,min=1,max=100"`
	Category string `json:"category" form:"category" validate:"required,min=1,max=50"`
	Desc     string `json:"desc" form:"desc" validate:"min=0,max=100"`
}

// 更新接口结构体
type UpdateApiRequest struct {
	Method   string `json:"Method" form:"Method" validate:"min=1,max=20"`
	Path     string `json:"Path" form:"Path" validate:"min=1,max=100"`
	Category string `json:"Category" form:"Category" validate:"min=1,max=50"`
	Desc     string `json:"Desc" form:"Desc" validate:"min=0,max=100"`
}

// 批量删除接口结构体
type DeleteApiRequest struct {
	ApiIds []uint `json:"ApiIds" form:"ApiIds"`
}

type ApiTreeDto struct {
	ID       int          `json:"ID"`
	Desc     string       `json:"Desc"`
	Category string       `json:"Category"`
	Children []*model.Api `json:"Children"`
}
