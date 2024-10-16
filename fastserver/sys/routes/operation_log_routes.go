package routes

import (
	"fastgin/sys/controller"
	"github.com/gin-gonic/gin"
)

func InitOperationLogRoutes(r *gin.RouterGroup) gin.IRoutes {
	operationLogController := controller.NewOperationLogController()
	router := r.Group("/log")
	{
		router.GET("/index", operationLogController.List)
		router.DELETE("/index", operationLogController.BatchDeleteByIds)
	}
	return r
}
