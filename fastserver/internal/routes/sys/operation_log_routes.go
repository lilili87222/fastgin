package sys

import (
	"fastgin/internal/controller/sys"
	"github.com/gin-gonic/gin"
)

func InitOperationLogRoutes(r *gin.RouterGroup) gin.IRoutes {
	operationLogController := sys.NewOperationLogController()
	router := r.Group("/log")
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/operation/list", operationLogController.GetOperationLogs)
		router.DELETE("/operation/delete/batch", operationLogController.BatchDeleteOperationLogByIds)
	}
	return r
}
