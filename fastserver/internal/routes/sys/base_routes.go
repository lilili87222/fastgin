package sys

import (
	"fastgin/internal/middleware"
	"github.com/gin-gonic/gin"
)

func InitBaseRoutes(router *gin.RouterGroup) gin.IRoutes {
	authMiddleware := middleware.GetJwtMiddleware()
	router.POST("/login", authMiddleware.LoginHandler)

	router.POST("/logout", authMiddleware.LogoutHandler)

	router.POST("/refreshToken", authMiddleware.RefreshHandler)

	return router
}