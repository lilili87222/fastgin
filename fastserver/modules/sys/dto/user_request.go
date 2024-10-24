package dto

// 用户登录结构体
type LoginRequest struct {
	UserName    string  `form:"user_name" json:"user_name" binding:"required"`
	Password    string  `form:"password" json:"password" binding:"required"`
	CaptchaId   string  `form:"captcha_id" json:"captcha_id" binding:"required"`
	CaptchaCode float64 `form:"captcha_code" json:"captcha_code" binding:"required"`
}

// 创建用户结构体
type CreateUserRequest struct {
	UserName string   `form:"user_name" json:"user_name" validate:"required,min=2,max=20"`
	Password string   `form:"password" json:"password"`
	Mobile   string   `form:"mobile" json:"mobile" validate:"required,checkMobile"`
	Avatar   string   `form:"avatar" json:"avatar"`
	NickName string   `form:"nick_name" json:"nick_name" validate:"min=0,max=20"`
	Des      string   `form:"des" json:"des" validate:"min=0,max=255"`
	Status   uint     `form:"status" json:"status" validate:"oneof=1 2"`
	RoleIds  []uint64 `form:"role_ids" json:"role_ids" validate:"required"`
}

// 获取用户列表结构体
type UserListRequest struct {
	Username string `json:"user_name" form:"user_name" `
	Mobile   string `json:"mobile" form:"mobile" `
	NickName string `json:"nick_name" form:"nick_name" `
	Status   uint   `json:"status" form:"status" `
	PageNum  uint   `json:"page_num" form:"page_num"`
	PageSize uint   `json:"page_size" form:"page_size"`
}

// 更新密码结构体
type ChangePwdRequest struct {
	OldPassword string `json:"old_password" form:"old_password" validate:"required"`
	NewPassword string `json:"new_password" form:"new_password" validate:"required"`
}

// 注册用户和重置密码
type RegisterRequest struct {
	UserName     string `form:"user_name" json:"user_name" binding:"required"`
	Password     string `form:"password" json:"password" binding:"required"`
	Repassword   string `form:"repassword" json:"repassword" binding:"required"`
	VerifyCodeId string `form:"verify_code_id" json:"verify_code_id" binding:"required"`
	VerifyCode   string `form:"verify_code" json:"verify_code" binding:"required"`
	Action       string `form:"action" json:"action" binding:"required"`
}
type SendVerifyCodeRequest struct {
	UserName    string  `form:"user_name" json:"user_name" binding:"required"`
	CaptchaId   string  `form:"captcha_id" json:"captcha_id" binding:"required"`
	CaptchaCode float64 `form:"captcha_code" json:"captcha_code" binding:"required"`
}
