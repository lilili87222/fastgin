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
	"github.com/google/uuid"
	"math"
	"math/rand"
	"time"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
//var store = base64Captcha.DefaultMemStore

type BaseController struct{}

const AngleSpin = 20 //角度偏差
// 获取验证码
// @Summary 获取验证码
// @Description 获取验证码
// @Tags 公开接口
// @Accept json
// @Produce json
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/public/captcha [get]
func (b *BaseController) Captcha(c *gin.Context) {
	angle := AngleSpin + rand.Float64()*(340-AngleSpin) // Generates a random float number between 20 and 340
	filePath := "conf/captcha/2.png"
	base64Img, err := util.Base64ImageFile(filePath, angle)
	if err != nil {
		httpz.BadRequest(c, "验证码获取失败"+err.Error())
		return
	}
	udidString := uuid.New().String()
	cache.Cache.Set(udidString, angle, 180*time.Second)
	httpz.Success(c, gin.H{"captchaId": udidString, "image": base64Img})
}

// 注册用户
// @Summary 注册用户
// @Description 注册用户
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param RegisterRequest body dto.RegisterRequest true "Register user request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/public/register [post]
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
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param SendVerifyCodeRequest body dto.SendVerifyCodeRequest true "Send verify code request"
// @Success 200 {object} httpz.ResponseBody
// @Failure 400 {object} httpz.ResponseBody
// @Router /api/public/verifycode [get]
func (b *BaseController) SendVerifyCode(c *gin.Context) {
	var req dto.SendVerifyCodeRequest
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	cid, found := cache.Cache.Get(req.CaptchaId)
	if !found {
		httpz.BadRequest(c, "验证码过期")
		return
	}
	angle := cid.(float64)
	if math.Abs(angle-req.CaptchaCode) > AngleSpin {
		httpz.BadRequest(c, "验证码错误")
		return
	}
	code := util.RandomString(6)
	codeID := uuid.New().String()

	e := email.SendRegisterEmail(req.UserName, "kxapp 验证码", code)
	if e != nil {
		config.Log.Error("验证码发送失败:" + e.Error())
		httpz.BadRequest(c, "验证码发送失败:"+e.Error())
		return
	}
	cache.Cache.Set(codeID, code, 1800*time.Second)
	httpz.Success(c, gin.H{"verify_code_id": codeID, "captchaLength": 6})
}
