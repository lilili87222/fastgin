package routes

import (
	"fastgin/sys/controller"
	"github.com/gin-gonic/gin"
)

func InitApiRoutes(r *gin.RouterGroup) gin.IRoutes {
	apiController := controller.NewApiController()
	router := r.Group("/api")
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", apiController.GetApis)
		router.GET("/tree", apiController.GetApiTree)
		router.POST("/create", apiController.CreateApi)
		router.PATCH("/update/:apiId", apiController.UpdateApiById)
		router.DELETE("/delete/batch", apiController.BatchDeleteApiByIds)
	}

	return r
}
