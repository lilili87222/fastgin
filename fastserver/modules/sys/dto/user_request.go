package dto

// 用户登录结构体
type RegisterAndLoginRequest struct {
	UserName string `form:"UserName" json:"UserName" binding:"required"`
	Password string `form:"Password" json:"Password" binding:"required"`
}

// 创建用户结构体
type CreateUserRequest struct {
	Username     string `form:"UserName" json:"UserName" validate:"required,min=2,max=20"`
	Password     string `form:"Password" json:"Password"`
	Mobile       string `form:"Mobile" json:"Mobile" validate:"required,checkMobile"`
	Avatar       string `form:"Avatar" json:"Avatar"`
	NickName     string `form:"NickName" json:"NickName" validate:"min=0,max=20"`
	Introduction string `form:"Introduction" json:"Introduction" validate:"min=0,max=255"`
	Status       uint   `form:"Status" json:"Status" validate:"oneof=1 2"`
	RoleIds      []uint `form:"RoleIds" json:"RoleIds" validate:"required"`
}

// 获取用户列表结构体
type UserListRequest struct {
	Username string `json:"UserName" form:"UserName" `
	Mobile   string `json:"Mobile" form:"Mobile" `
	NickName string `json:"NickName" form:"NickName" `
	Status   uint   `json:"Status" form:"Status" `
	PageNum  uint   `json:"PageNum" form:"PageNum"`
	PageSize uint   `json:"PageSize" form:"PageSize"`
}

// 批量删除用户结构体
//type DeleteUserRequest struct {
//	UserIds []uint `json:"UserIds" form:"UserIds"`
//}

// 更新密码结构体
type ChangePwdRequest struct {
	OldPassword string `json:"OldPassword" form:"OldPassword" validate:"required"`
	NewPassword string `json:"NewPassword" form:"NewPassword" validate:"required"`
}
