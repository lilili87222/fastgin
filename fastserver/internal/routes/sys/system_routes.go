package sys

import (
	"fastgin/internal/controller/sys"
	"github.com/gin-gonic/gin"
)

func InitSystemRoutes(r *gin.RouterGroup) gin.IRoutes {
	menuController := sys.NewSystemController()
	router := r.Group("/system")
	{
		router.GET("/info", menuController.GetSystemInformation)
		router.GET("/stop", menuController.GetStopServer)
		router.GET("/restart", menuController.RestartServer)
	}
	return r
}
