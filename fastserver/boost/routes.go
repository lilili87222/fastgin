package boost

import (
	"fastgin/config"
	_ "fastgin/docs" // Import the generated docs
	"fastgin/sys/middleware"
	"fastgin/sys/routes"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

// 初始化
func InitRoutes(engine *gin.Engine) {

	// 创建不带中间件的路由:
	// engine := gin.New()
	// engine.Use(gin.Recovery())

	// 启用限流中间件
	// 默认每50毫秒填充一个令牌，最多填充200个

	// 启用全局跨域中间件
	engine.Use(middleware.CORSMiddleware())

	engine.Group("/").GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	publicGroup := engine.Group("api/public")
	routes.InitBaseRoutes(publicGroup) // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件

	authGroup := engine.Group(config.Instance.System.UrlPathPrefix)
	// 启用中间件
	authGroup.Use(middleware.OperationLogMiddleware())
	authGroup.Use(middleware.RateLimitMiddleware(time.Millisecond*time.Duration(config.Instance.RateLimit.FillInterval), config.Instance.RateLimit.Capacity))
	authGroup.Use(middleware.GetJwtMiddleware().MiddlewareFunc()) // jwt认证中间件
	authGroup.Use(middleware.CasbinMiddleware())                  //// 开启casbin鉴权中间件

	routes.InitUserRoutes(authGroup)         // 注册用户路由, jwt认证中间件,casbin鉴权中间件
	routes.InitRoleRoutes(authGroup)         // 注册角色路由, jwt认证中间件,casbin鉴权中间件
	routes.InitMenuRoutes(authGroup)         // 注册菜单路由, jwt认证中间件,casbin鉴权中间件
	routes.InitApiRoutes(authGroup)          // 注册接口路由, jwt认证中间件,casbin鉴权中间件
	routes.InitOperationLogRoutes(authGroup) // 注册操作日志路由, jwt认证中间件,casbin鉴权中间件

	routes.InitSystemRoutes(authGroup) // 注册系统路由, jwt认证中间件,casbin鉴权中间件

	config.Log.Info("初始化路由完成！")
}