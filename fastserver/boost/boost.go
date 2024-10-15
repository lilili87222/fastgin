package boost

import (
	"context"
	"errors"
	"fastgin/config"
	"fastgin/sys/middleware"
	"fastgin/sys/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var httpServer *http.Server

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
	InitRoutes(engine)

	httpServer = &http.Server{Addr: fmt.Sprintf(":" + config.Instance.System.Port), Handler: engine}

	// 启动服务器的 goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	if err := httpServer.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %v\n", err)
	}
	fmt.Println("Server exiting")
}
