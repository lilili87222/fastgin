package dto

type UsersDto struct {
	Id       uint64   `json:"id"`
	UserName string   `json:"user_name"`
	Mobile   string   `json:"mobile"`
	Avatar   string   `json:"avatar"`
	NickName string   `json:"nick_name"`
	Des      string   `json:"des"`
	Status   uint     `json:"status"`
	Creator  string   `json:"creator"`
	RoleIds  []uint64 `json:"role_ids"`
}

//type UserInfoDto struct {
//	Id           uint64        `json:"ID"`
//	UserName     string        `json:"UserName"`
//	Mobile       string        `json:"Mobile"`
//	Avatar       string        `json:"Avatar"`
//	NickName     string        `json:"NickName"`
//	Des string        `json:"Des"`
//	Roles        []*model.Role `json:"Roles"`
//}
