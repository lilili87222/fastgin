package main

import (
	"errors"
	"fastgin/config"
	"fastgin/internal/dao/sys"
	"fastgin/internal/middleware"
	"fastgin/internal/routes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @title Go Web fastgin API
// @version 1.0
// @description This is a sample server for a Go web mini project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 192.168.123.214:8088
// @BasePath /
func main() {

	// 加载配置文件到全局配置结构体
	config.InitConfig()

	// 初始化日志
	config.InitLogger()

	// 初始化数据库
	config.InitDatabase()
	// 初始化casbin策略管理器
	config.InitCasbinEnforcer()
	// 初始化数据
	config.InitData()

	// 初始化Validator数据校验
	config.InitValidate()

	// 操作日志中间件处理日志时没有将日志发送到rabbitmq或者kafka中, 而是发送到了channel中
	// 这里开启3个goroutine处理channel将日志记录到数据库
	logRepository := sys.NewOperationLogRepository()
	for i := 0; i < 3; i++ {
		go logRepository.SaveOperationLogChannel(middleware.OperationLogChan)
	}

	//设置模式
	gin.SetMode(config.Conf.System.Mode)
	engine := gin.Default()
	routes.InitRoutes(engine)

	server := &http.Server{Addr: fmt.Sprintf(":" + config.Conf.System.Port), Handler: engine}
	if err := server.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
		config.Log.Fatalf("listen: %s\n", err)
	}
}
