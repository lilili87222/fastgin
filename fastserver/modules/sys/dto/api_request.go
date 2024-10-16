package dto

import "fastgin/modules/sys/model"

type ApiTreeDto struct {
	Id       int          `json:"Id"`
	Desc     string       `json:"Desc"`
	Category string       `json:"Category"`
	Children []*model.Api `json:"Children"`
}

// 获取接口列表结构体
type ApiListRequest struct {
	Method   string `json:"Method" form:"Method"`
	Path     string `json:"Path" form:"Path"`
	Category string `json:"Category" form:"Category"`
	Creator  string `json:"Creator" form:"Creator"`
	PageNum  uint   `json:"PageNum" form:"PageNum"`
	PageSize uint   `json:"PageSize" form:"PageSize"`
}

// 创建接口结构体
type CreateApiRequest struct {
	Method   string `json:"Method" form:"Method" validate:"required,min=1,max=20"`
	Path     string `json:"Path" form:"Path" validate:"required,min=1,max=100"`
	Category string `json:"Category" form:"Category" validate:"required,min=1,max=50"`
	Desc     string `json:"Desc" form:"Desc" validate:"min=0,max=100"`
}
