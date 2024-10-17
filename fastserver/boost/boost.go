package boost

import (
	"context"
	"errors"
	"fastgin/config"
	_ "fastgin/docs" // Import the generated docs
	appRoute "fastgin/modules/app/route"
	"fastgin/modules/sys/middleware"
	"fastgin/modules/sys/route"
	"fastgin/modules/sys/script"
	"fastgin/modules/sys/service"
	"fmt"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var HttpServer *http.Server
var Engine *gin.Engine
var AuthGroup *gin.RouterGroup
var PublicGroup *gin.RouterGroup

func StartWebService() {
	// 操作日志中间件处理日志时没有将日志发送到rabbitmq或者kafka中, 而是发送到了channel中
	// 这里开启3个goroutine处理channel将日志记录到数据库
	logDao := service.NewLogService()
	for i := 0; i < 3; i++ {
		go logDao.SaveOperationLogChannel(middleware.OperationLogChan)
	}
	//设置模式
	gin.SetMode(config.Instance.System.Mode)
	engine := gin.Default()
	initRoutes(engine)

	HttpServer = &http.Server{Addr: fmt.Sprintf(":" + config.Instance.System.Port), Handler: engine}

	// 启动服务器的 goroutine
	go func() {
		if err := HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("Listen: %s\n", err)
		}
	}()
	// 捕获信号（可选）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit // 等待信号

	// 优雅地停止服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := HttpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}
	fmt.Println("Server exiting")
}

// 初始化
func initRoutes(engine *gin.Engine) {

	// 创建不带中间件的路由:
	// engine := gin.New()
	// engine.Use(gin.Recovery())

	// 启用限流中间件
	// 默认每50毫秒填充一个令牌，最多填充200个

	// 启用全局跨域中间件
	engine.Use(middleware.CORSMiddleware())
	engine.Group("/").GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	PublicGroup = engine.Group("api/public")
	AuthGroup = engine.Group(config.Instance.System.UrlPathPrefix)
	AuthGroup.Use(middleware.OperationLogMiddleware())
	AuthGroup.Use(middleware.RateLimitMiddleware(time.Millisecond*time.Duration(config.Instance.RateLimit.FillInterval), config.Instance.RateLimit.Capacity))
	AuthGroup.Use(middleware.GetJwtMiddleware().MiddlewareFunc()) // jwt认证中间件
	AuthGroup.Use(middleware.CasbinMiddleware())                  //// 开启casbin鉴权中间件

	script.InitSysModuleDatabase()
	initApiRoutes()

	config.Log.Info("初始化路由完成！")
}

func initApiRoutes() {
	route.InitBaseRoutes(PublicGroup)       // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件
	route.InitUserRoutes(AuthGroup)         // 注册用户路由, jwt认证中间件,casbin鉴权中间件
	route.InitRoleRoutes(AuthGroup)         // 注册角色路由, jwt认证中间件,casbin鉴权中间件
	route.InitMenuRoutes(AuthGroup)         // 注册菜单路由, jwt认证中间件,casbin鉴权中间件
	route.InitApiRoutes(AuthGroup)          // 注册接口路由, jwt认证中间件,casbin鉴权中间件
	route.InitOperationLogRoutes(AuthGroup) // 注册操作日志路由, jwt认证中间件,casbin鉴权中间件
	route.InitSystemRoutes(AuthGroup)       // 注册系统路由, jwt认证中间件,casbin鉴权中间件

	appRoute.InitDictionary(AuthGroup)
}
