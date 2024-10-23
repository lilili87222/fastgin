package controller

import (
	"fastgin/boost/config"
	"fastgin/common/cache"
	"fastgin/common/email"
	"fastgin/common/httpz"
	"fastgin/common/util"
	"fastgin/modules/sys/dto"
	"fastgin/modules/sys/model"
	"fastgin/modules/sys/service"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
var store = base64Captcha.DefaultMemStore

type BaseController struct{}

// 获取验证码
// @Summary 获取验证码
// @Description 获取验证码
// @Tags Base
// @Accept json
// @Produce json
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/captcha [get]
func (b *BaseController) Captcha(c *gin.Context) {
	capConfig := config.Configs.Captcha
	driver := base64Captcha.NewDriverDigit(capConfig.ImgHeight, capConfig.ImgWidth, capConfig.KeyLong, 0.7, 80)
	cp := base64Captcha.NewCaptcha(driver, store)
	id, b64s, _, err := cp.Generate()
	if err != nil {
		httpz.BadRequest(c, "验证码获取失败")
		return
	}
	httpz.Success(c, gin.H{"captchaId": id, "image": b64s, "captchaLength": capConfig.KeyLong})
}

// 注册用户
// @Summary 注册用户
// @Description 注册用户
// @Tags Base
// @Accept json
// @Produce json
// @Param RegisterRequest body dto.RegisterRequest true "Register user request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/register [post]
func (b *BaseController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	if req.Password != req.Repassword {
		httpz.BadRequest(c, "两次密码不一致")
		return
	}
	if len(req.Password) < 8 {
		httpz.BadRequest(c, "密码长度不能小于8位")
		return
	}
	cacheCode := cache.GetString(req.VerifyCodeId)
	if cacheCode == "" || cacheCode != req.VerifyCode {
		httpz.BadRequest(c, "验证码错误")
		return
	}
	userService := service.NewUserService()
	u, _ := userService.GetUserByUsername(req.UserName)
	//registor:=u==nil
	if u == nil {
		user := model.User{
			UserName: req.UserName,
			Password: req.Password,
			Status:   1,
			Creator:  "register",
			Roles:    []model.Role{{ID: 2}},
		}
		if util.IsPhoneNumber(req.UserName) {
			user.Mobile = req.UserName
		}
		err := userService.CreateUser(&user)
		if err != nil {
			httpz.ServerError(c, "注册失败: "+err.Error())
			return
		}
	} else {
		err := userService.ChangePwd(req.UserName, req.Password)
		if err != nil {
			httpz.ServerError(c, "更新密码失败: "+err.Error())
			return
		}
	}
	httpz.Success(c, nil)

}

// 发送验证码
// @Summary 发送验证码
// @Description 发送验证码
// @Tags Base
// @Accept json
// @Produce json
// @Param SendVerifyCodeRequest body dto.SendVerifyCodeRequest true "Send verify code request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/auth/verifycode [get]
func (b *BaseController) SendVerifyCode(c *gin.Context) {
	var req dto.SendVerifyCodeRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	if !store.Verify(req.CaptchaId, req.CaptchaCode, true) {
		httpz.BadRequest(c, "验证码错误")
		return
	}
	code := util.RandomString(6)
	codeID := util.RandomString(24)

	e := email.SendRegisterEmail(req.UserName, "kxapp 验证码", code)
	if e != nil {
		config.Log.Error("验证码发送失败:" + e.Error())
		httpz.BadRequest(c, "验证码发送失败:"+e.Error())
		return
	}
	cache.SetString(codeID, code)
	httpz.Success(c, gin.H{"verify_code_id": codeID, "captchaLength": 6})
}
