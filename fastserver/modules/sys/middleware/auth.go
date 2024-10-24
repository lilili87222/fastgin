package middleware

import (
	config2 "fastgin/boost/config"
	"fastgin/common/cache"
	"fastgin/common/httpz"
	"fastgin/common/util"
	"fastgin/modules/sys/dto"
	"fastgin/modules/sys/middleware/jwt"
	"fastgin/modules/sys/model"
	"fastgin/modules/sys/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"time"
)

var jwtMiddleware *jwt.GinJWTMiddleware

// 初始化jwt中间件
func GetJwtMiddleware() *jwt.GinJWTMiddleware {
	if jwtMiddleware != nil {
		return jwtMiddleware
	}
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:           config2.Configs.Jwt.Realm,                                 // jwt标识
		Key:             []byte(config2.Configs.Jwt.Key),                           // 服务端密钥
		Timeout:         time.Hour * time.Duration(config2.Configs.Jwt.Timeout),    // token过期时间
		MaxRefresh:      time.Hour * time.Duration(config2.Configs.Jwt.MaxRefresh), // token最大刷新时间(RefreshToken过期时间=Timeout+MaxRefresh)
		PayloadFunc:     payloadFunc,                                               // 有效载荷处理
		IdentityHandler: identityHandler,                                           // 解析Claims
		Authenticator:   login,                                                     // 校验token的正确性, 处理登录逻辑
		Authorizator:    authorizator,                                              // 用户登录校验成功处理
		Unauthorized:    unauthorized,                                              // 用户登录校验失败处理
		LoginResponse:   loginResponse,                                             // 登录成功后的响应
		LogoutResponse:  logoutResponse,                                            // 登出后的响应
		RefreshResponse: refreshResponse,                                           // 刷新token后的响应
		TokenLookup:     "header: Authorization, query: token, cookie: jwt",        // 自动在这几个地方寻找请求中的token
		TokenHeadName:   "Bearer",                                                  // header名称
		TimeFunc:        time.Now,
	})
	if err != nil {
		config2.Log.Panicf("初始化JWT中间件失败：%v", err)
		panic(fmt.Sprintf("初始化JWT中间件失败：%v", err))
	}
	jwtMiddleware = authMiddleware
	return authMiddleware
}

// 有效载荷处理
func payloadFunc(data interface{}) jwt.MapClaims {
	if v, ok := data.(map[string]interface{}); ok {
		var user model.User
		// 将用户json转为结构体
		util.JsonI2Struct(v["user"], &user)
		return jwt.MapClaims{
			jwt.IdentityKey: user.ID,
			"user":          v["user"],
		}
	}
	return jwt.MapClaims{}
}

// 解析Claims
func identityHandler(c *gin.Context) interface{} {
	claims := jwt.ExtractClaims(c)
	// 此处返回值类型map[string]interface{}与payloadFunc和authorizator的data类型必须一致, 否则会导致授权失败还不容易找到原因
	return map[string]interface{}{
		"IdentityKey": claims[jwt.IdentityKey],
		"user":        claims["user"],
	}
}

// 校验token的正确性, 处理登录逻辑
func login(c *gin.Context) (interface{}, error) {
	var req dto.LoginRequest
	// 请求json绑定
	if err := c.ShouldBind(&req); err != nil {
		return "", err
	}
	cid, found := cache.Cache.Get(req.CaptchaId)
	if !found {
		return nil, fmt.Errorf("验证码已过期")
	}
	angle := cid.(float64)
	if math.Abs(angle-req.CaptchaCode) > 20 {
		return nil, fmt.Errorf("验证码错误")
	}
	u := &model.User{
		UserName: req.UserName,
		Password: req.Password,
	}
	// 密码校验
	userService := service.NewUserService()
	user, err := userService.Login(u)
	if err != nil {
		return nil, err
	}
	// 将用户以json格式写入, payloadFunc/authorizator会使用到
	return map[string]interface{}{
		"user": util.Struct2Json(user),
	}, nil
}

// 用户登录校验成功处理
func authorizator(data interface{}, c *gin.Context) bool {
	if v, ok := data.(map[string]interface{}); ok {
		userStr := v["user"].(string)
		var user model.User
		// 将用户json转为结构体
		util.Json2Struct(userStr, &user)
		// 将用户保存到context, api调用时取数据方便
		c.Set("user", user)
		return true
	}
	return false
}

// 用户登录校验失败处理
func unauthorized(c *gin.Context, code int, message string) {
	config2.Log.Debugf("JWT认证失败, 错误码: %d, 错误信息: %s", code, message)
	httpz.Response(c, code, nil, fmt.Sprintf("JWT认证失败, 错误码: %d, 错误信息: %s", code, message))
}

// 登录成功后的响应
func loginResponse(c *gin.Context, code int, token string, expires time.Time) {
	httpz.Response(c, code,
		gin.H{
			"token":   token,
			"expires": expires.Format("2006-01-02 15:04:05"),
		},
		"登录成功")
}

// 登出后的响应
func logoutResponse(c *gin.Context, code int) {
	httpz.Success(c, nil)
}

// 刷新token后的响应
func refreshResponse(c *gin.Context, code int, token string, expires time.Time) {
	httpz.Response(c, code,
		gin.H{
			"token":   token,
			"expires": expires,
		},
		"刷新token成功")
}
