package boost

import (
	"context"
	"errors"
	"fastgin/boost/config"
	"fastgin/common/storage"
	"fastgin/database"
	_ "fastgin/docs" // Import the generated docs
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

func preInit() {
	// 加载配置文件到全局配置结构体
	config.InitConfig()

	// 初始化日志
	config.InitLogger()

	// 初始化数据库
	database.InitDatabaseConnection()
	// 初始化casbin策略管理器
	config.InitCasbinEnforcer(database.DB)

	// 创建数据库表，插入初始数据
	script.InitSysModuleDatabase()

	// 初始化Validator数据校验
	config.InitValidate()

	storage.InitStorage(config.Configs.Storage)
}
func StartWebService() {
	preInit()
	// 操作日志中间件处理日志时没有将日志发送到rabbitmq或者kafka中, 而是发送到了channel中
	// 这里开启3个goroutine处理channel将日志记录到数据库
	logService := service.NewLogService()
	for i := 0; i < 3; i++ {
		go logService.SaveOperationLogChannel(middleware.OperationLogChan)
	}
	//设置模式
	gin.SetMode(config.Configs.System.Mode)
	engine := gin.Default()
	p, a := initMiddlewares(engine)
	initApiRoutes(engine, p, a)

	HttpServer = &http.Server{Addr: fmt.Sprintf(":" + config.Configs.System.Port), Handler: engine}

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
func initMiddlewares(engine *gin.Engine) (*gin.RouterGroup, *gin.RouterGroup) {

	// 创建不带中间件的路由:
	// engine := gin.New()
	// engine.Use(gin.Recovery())

	// 启用限流中间件
	// 默认每50毫秒填充一个令牌，最多填充200个

	// 启用全局跨域中间件
	engine.Use(middleware.CORSMiddleware())
	engine.Group("/").GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	PublicGroup := engine.Group("api/public")
	AuthGroup := engine.Group(config.Configs.System.UrlPathPrefix)
	AuthGroup.Use(middleware.OperationLogMiddleware())
	AuthGroup.Use(middleware.RateLimitMiddleware(time.Millisecond*time.Duration(config.Configs.RateLimit.FillInterval), config.Configs.RateLimit.Capacity))
	AuthGroup.Use(middleware.GetJwtMiddleware().MiddlewareFunc()) // jwt认证中间件
	AuthGroup.Use(middleware.CasbinMiddleware())                  //// 开启casbin鉴权中间件
	config.Log.Info("初始化Middleware完成！")
	return PublicGroup, AuthGroup
}

func initApiRoutes(engine *gin.Engine, PublicGroup *gin.RouterGroup, AuthGroup *gin.RouterGroup) {
	route.InitBaseRoutes(PublicGroup)       // 注册基础路由, 不需要jwt认证中间件,不需要casbin中间件
	route.InitUserRoutes(AuthGroup)         // 注册用户路由, jwt认证中间件,casbin鉴权中间件
	route.InitRoleRoutes(AuthGroup)         // 注册角色路由, jwt认证中间件,casbin鉴权中间件
	route.InitMenuRoutes(AuthGroup)         // 注册菜单路由, jwt认证中间件,casbin鉴权中间件
	route.InitApiRoutes(AuthGroup)          // 注册接口路由, jwt认证中间件,casbin鉴权中间件
	route.InitOperationLogRoutes(AuthGroup) // 注册操作日志路由, jwt认证中间件,casbin鉴权中间件
	route.InitSystemRoutes(AuthGroup)       // 注册系统路由, jwt认证中间件,casbin鉴权中间件
	route.InitDictionary(AuthGroup)
	config.Log.Info("初始化路由完成！")
}
