package sys

import (
	"go-web-mini/internal/model/sys"
)

// 返回给前端的当前用户信息
type UserInfoDto2 struct {
	ID           uint        `json:"id"`
	Username     string      `json:"username"`
	Mobile       string      `json:"mobile"`
	Avatar       string      `json:"avatar"`
	Nickname     string      `json:"nickname"`
	Introduction string      `json:"introduction"`
	Roles        []*sys.Role `json:"roles"`
}

func ToUserInfoDto(user sys.User) UserInfoDto2 {
	return UserInfoDto2{
		ID:           user.ID,
		Username:     user.Username,
		Mobile:       user.Mobile,
		Avatar:       user.Avatar,
		Nickname:     *user.Nickname,
		Introduction: *user.Introduction,
		Roles:        user.Roles,
	}
}

// 返回给前端的用户列表
type UsersDto2 struct {
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

func ToUsersDto(userList []*sys.User) []UsersDto2 {
	var users []UsersDto2
	for _, user := range userList {
		userDto := UsersDto2{
			ID:           user.ID,
			Username:     user.Username,
			Mobile:       user.Mobile,
			Avatar:       user.Avatar,
			Nickname:     *user.Nickname,
			Introduction: *user.Introduction,
			Status:       user.Status,
			Creator:      user.Creator,
		}
		roleIds := make([]uint, 0)
		for _, role := range user.Roles {
			roleIds = append(roleIds, role.ID)
		}
		userDto.RoleIds = roleIds
		users = append(users, userDto)
	}

	return users
}
