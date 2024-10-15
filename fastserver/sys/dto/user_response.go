package dto

import "fastgin/sys/model"

type UsersDto struct {
	ID           uint   `json:"ID"`
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
	ID           uint          `json:"ID"`
	UserName     string        `json:"UserName"`
	Mobile       string        `json:"Mobile"`
	Avatar       string        `json:"Avatar"`
	NickName     string        `json:"NickName"`
	Introduction string        `json:"Introduction"`
	Roles        []*model.Role `json:"Roles"`
}
