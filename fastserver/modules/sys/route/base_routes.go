package route

import (
	"fastgin/modules/sys/middleware"
	"github.com/gin-gonic/gin"
)

func InitBaseRoutes(router *gin.RouterGroup) gin.IRoutes {
	authMiddleware := middleware.GetJwtMiddleware()
	router.POST("/login", authMiddleware.LoginHandler)
	return router
}
