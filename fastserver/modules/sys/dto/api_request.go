package dto

import "fastgin/modules/sys/model"

type ApiTreeDto struct {
	Id       int          `json:"id"`
	Des      string       `json:"des"`
	Category string       `json:"category"`
	Children []*model.Api `json:"children"`
}

// 获取接口列表结构体
type ApiListRequest struct {
	Method   string `json:"method" form:"method"`
	Path     string `json:"path" form:"path"`
	Category string `json:"category" form:"category"`
	Creator  string `json:"creator" form:"creator"`
	PageNum  uint   `json:"page_num" form:"page_num"`
	PageSize uint   `json:"page_size" form:"page_size"`
}

// 创建接口结构体
type CreateApiRequest struct {
	Method   string `json:"method" form:"method" validate:"required,min=1,max=20"`
	Path     string `json:"path" form:"path" validate:"required,min=1,max=100"`
	Category string `json:"category" form:"category" validate:"required,min=1,max=50"`
	Des      string `json:"des" form:"des" validate:"min=0,max=100"`
}
