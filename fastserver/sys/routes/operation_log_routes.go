package routes

import (
	"fastgin/sys/controller"
	"github.com/gin-gonic/gin"
)

func InitOperationLogRoutes(r *gin.RouterGroup) gin.IRoutes {
	operationLogController := controller.NewOperationLogController()
	router := r.Group("/log")
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/operation/list", operationLogController.GetOperationLogs)
		router.DELETE("/operation/delete/batch", operationLogController.BatchDeleteOperationLogByIds)
	}
	return r
}
