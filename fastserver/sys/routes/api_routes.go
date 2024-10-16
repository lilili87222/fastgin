package routes

import (
	"fastgin/sys/controller"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup) gin.IRoutes {
	apiController := controller.NewApiController()
	router := r.Group("/api")
	{
		router.GET("/index", apiController.List)
		router.POST("/index", apiController.Create)
		router.PATCH("/index/:apiId", apiController.UpdateById)
		router.DELETE("/index", apiController.BatchDeleteByIds)

		router.GET("/tree", apiController.GetApiTree)
	}

	return r
}
