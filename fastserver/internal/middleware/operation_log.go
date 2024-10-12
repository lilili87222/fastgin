package middleware

import (
	"fastgin/config"
	sys2 "fastgin/internal/dao/sys"
	"fastgin/internal/model/sys"
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

// 操作日志channel
var OperationLogChan = make(chan *sys.OperationLog, 30)

func OperationLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行耗时
		timeCost := endTime.Sub(startTime).Milliseconds()

		// 获取当前登录用户
		var username string
		ctxUser, exists := c.Get("user")
		if !exists {
			username = "未登录"
		}
		user, ok := ctxUser.(sys.User)
		if !ok {
			username = "未登录"
		}
		username = user.Username

		//requestPath := c.Request.RequestURI
		fullPath := c.FullPath()
		requestURL := c.Request.RequestURI
		fmt.Println("requestPath:", requestURL, "fullPath:", fullPath)
		// 获取访问路径
		path := strings.TrimPrefix(fullPath, "/"+config.Conf.System.UrlPathPrefix)
		// 请求方式
		method := c.Request.Method

		// 获取接口描述
		apiRepository := sys2.NewApiRepository()
		apiDesc, _ := apiRepository.GetApiDescByPath(path, method)

		operationLog := sys.OperationLog{
			Username:   username,
			Ip:         c.ClientIP(),
			IpLocation: "",
			Method:     method,
			Path:       path,
			Desc:       apiDesc,
			Status:     c.Writer.Status(),
			StartTime:  startTime,
			TimeCost:   timeCost,
			//UserAgent:  c.Request.UserAgent(),
		}

		// 最好是将日志发送到rabbitmq或者kafka中
		// 这里是发送到channel中，开启3个goroutine处理
		OperationLogChan <- &operationLog
	}
}
