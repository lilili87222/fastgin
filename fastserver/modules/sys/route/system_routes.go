package route

import (
	"fastgin/modules/sys/controller"
	"github.com/gin-gonic/gin"
)

func InitSystemRoutes(r *gin.RouterGroup) gin.IRoutes {
	menuController := controller.NewSystemController()
	router := r.Group("/system")
	{
		router.GET("/info", menuController.GetSystemInformation)
		router.GET("/stop", menuController.GetStopServer)
		router.GET("/restart", menuController.RestartServer)
	}
	return r
}
