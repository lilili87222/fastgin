package controller

import (
	"fastgin/boost/config"
	"fastgin/common/cache"
	"fastgin/common/email"
	"fastgin/common/httpz"
	"fastgin/common/util"
	"fastgin/modules/sys/dto"
	"fastgin/modules/sys/middleware"
	"fastgin/modules/sys/model"
	"fastgin/modules/sys/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math/rand"
	"os"
	"time"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
//var store = base64Captcha.DefaultMemStore

type BaseController struct{}

// @Summary 用户登录
// @Description 用户登录获取JWT Token {"UserName": "testlog", "Password": "123456"}
// @Tags 公开接口
// @Accept json
// @Produce json
// @Param login body dto.LoginRequest true "登录信息" default({"username": "testlog", "password": "123456"})
// @Success 200 {object} map[string]interface{} "{"code":200,"token":"xxx","expire":"xxx"}"
// @Failure 401 {object} map[string]interface{} "{"code":401,"message":"Unauthorized"}"
// @Router /api/public/login [post]
// LoginHandler can be used by clients to get a jwt token.
// Payload needs to be json in the form of {"username": "USERNAME", "password": "PASSWORD"}.
// Reply will be of the form {"token": "TOKEN"}.
func (b *BaseController) Login(c *gin.Context) {
	var req dto.LoginRequest
	// 请求json绑定
	if err := c.ShouldBind(&req); err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	cid, found := cache.Cache.Get(req.CaptchaId)
	if !found {
		httpz.BadRequest(c, "验证码已过期")
		return
	}
	if !util.EqualCaptcha(cid.(float64), req.CaptchaCode) {
		httpz.BadRequest(c, "验证码错误")
		return
	}
	u := &model.User{
		UserName: req.UserName,
		Password: req.Password,
	}
	// 密码校验
	userService := service.NewUserService()
	user, err := userService.Login(u)
	if err != nil {
		httpz.BadRequest(c, err.Error())
		return
	}
	token, err := middleware.GenerateJWTToken(user.ID)
	if err != nil {
		httpz.ServerError(c, "生成token失败")
		return
	}
	cache.SetUser(user)
	httpz.Success(c, gin.H{"token": token})
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
	user := userService.GetUserByUsername(req.UserName)
	//registor:=u==nil
	if user == nil {
		user = &model.User{
			UserName: req.UserName,
			Password: req.Password,
			Status:   1,
			Creator:  "register",
			Roles:    []model.Role{{ID: 2}},
		}
		if util.IsPhoneNumber(req.UserName) {
			user.Mobile = req.UserName
		}
		err := userService.CreateUser(user)
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

	token, err := middleware.GenerateJWTToken(user.ID)
	if err != nil {
		httpz.ServerError(c, "生成token失败")
		return
	}
	cache.SetUser(user)
	httpz.Success(c, gin.H{"token": token})

	//httpz.Success(c, nil)

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
	if !util.EqualCaptcha(cid.(float64), req.CaptchaCode) {
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
	httpz.Success(c, gin.H{"verify_code_id": codeID, "code_length": 6})
}

var imageFileList []string

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
	if len(imageFileList) == 0 {
		imgEnt, _ := os.ReadDir("conf/captcha")
		for _, v := range imgEnt {
			imageFileList = append(imageFileList, v.Name())
		}
	}
	seed := rand.NewSource(time.Now().UnixNano())
	r := rand.New(seed)
	index := r.Intn(len(imageFileList))
	angle := util.RandCaptchaAngle() // Generates a random float number between 20 and 340
	filePath := fmt.Sprintf("conf/captcha/%s", imageFileList[index])

	base64Img, err := util.Base64ImageFile(filePath, angle)
	if err != nil {
		httpz.BadRequest(c, "验证码获取失败"+err.Error())
		return
	}
	udidString := uuid.New().String()
	cache.Cache.Set(udidString, angle, 180*time.Second)
	httpz.Success(c, gin.H{"captcha_id": udidString, "image": base64Img})
}
