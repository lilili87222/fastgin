package route

import (
	"fastgin/modules/sys/controller"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup) gin.IRoutes {
	apiController := controller.NewApiController()
	router := r.Group("/api")
	{
		router.GET("/index", apiController.List)
		router.POST("/index", apiController.Create)
		router.PATCH("/index/:apiId", apiController.Update)
		router.DELETE("/index", apiController.BatchDelete)

		router.GET("/tree", apiController.GetApiTree)
	}

	return r
}
