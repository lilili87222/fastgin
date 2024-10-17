package dto

import "fastgin/modules/sys/model"

type UsersDto struct {
	Id           uint   `json:"Id"`
	UserName     string `json:"UserName"`
	Mobile       string `json:"Mobile"`
	Avatar       string `json:"Avatar"`
	NickName     string `json:"NickName"`
	Introduction string `json:"Introduction"`
	Status       uint   `json:"Status"`
	Creator      string `json:"Creator"`
	RoleIds      []uint `json:"RoleIds"`
}
type UserInfoDto struct {
	Id           uint          `json:"Id"`
	UserName     string        `json:"UserName"`
	Mobile       string        `json:"Mobile"`
	Avatar       string        `json:"Avatar"`
	NickName     string        `json:"NickName"`
	Introduction string        `json:"Introduction"`
	Roles        []*model.Role `json:"Roles"`
}
