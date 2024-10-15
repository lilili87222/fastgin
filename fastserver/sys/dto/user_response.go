package dto

type UsersDto struct {
	ID           uint   `json:"ID"`
	Username     string `json:"username"`
	Mobile       string `json:"mobile"`
	Avatar       string `json:"avatar"`
	Nickname     string `json:"nickname"`
	Introduction string `json:"introduction"`
	Status       uint   `json:"status"`
	Creator      string `json:"creator"`
	RoleIds      []uint `json:"roleIds"`
}
