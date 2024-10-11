package routes

import (
	"fastgin/config"
	_ "fastgin/docs" // Import the generated docs
	middleware2 "fastgin/internal/middleware"
	"fastgin/internal/routes/sys"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

// 初始化
func InitRoutes() *gin.Engine {
	//设置模式
	gin.SetMode(config.Conf.System.Mode)

	// 创建带有默认中间件的路由:
	// 日志与恢复中间件
	r := gin.Default()
	// 创建不带中间件的路由:
	// r := gin.New()
	// r.Use(gin.Recovery())

	// 启用限流中间件
	// 默认每50毫秒填充一个令牌，最多填充200个

	// 启用全局跨域中间件
	r.Use(middleware2.CORSMiddleware())

	r.Group("/").GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 路由分组
	apiGroup := r.Group("/" + config.Conf.System.UrlPathPrefix)
	// 启用操作日志中间件
	apiGroup.Use(middleware2.OperationLogMiddleware())
	apiGroup.Use(middleware2.RateLimitMiddleware(time.Millisecond*time.Duration(config.Conf.RateLimit.FillInterval), config.Conf.RateLimit.Capacity))

	// 初始化JWT认证中间件
	authMiddleware, err := middleware2.InitAuth()
	if err != nil {
		config.Log.Panicf("初始化JWT中间件失败：%v", err)
		panic(fmt.Sprintf("初始化JWT中间件失败：%v", err))
	}
	// 注册路由
	sys.InitBaseRoutes(apiGroup, authMiddleware)         // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件
	sys.InitUserRoutes(apiGroup, authMiddleware)         // 注册用户路由, jwt认证中间件,casbin鉴权中间件
	sys.InitRoleRoutes(apiGroup, authMiddleware)         // 注册角色路由, jwt认证中间件,casbin鉴权中间件
	sys.InitMenuRoutes(apiGroup, authMiddleware)         // 注册菜单路由, jwt认证中间件,casbin鉴权中间件
	sys.InitApiRoutes(apiGroup, authMiddleware)          // 注册接口路由, jwt认证中间件,casbin鉴权中间件
	sys.InitOperationLogRoutes(apiGroup, authMiddleware) // 注册操作日志路由, jwt认证中间件,casbin鉴权中间件

	config.Log.Info("初始化路由完成！")
	return r
}
