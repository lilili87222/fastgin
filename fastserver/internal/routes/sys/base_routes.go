package sys

import (
	"fastgin/internal/middleware"
	"github.com/gin-gonic/gin"
)

// 注册基础路由
func InitBaseRoutes(router *gin.RouterGroup) gin.IRoutes {
	//router := r.Group("/base")
	//{
	authMiddleware := middleware.GetJwtMiddleware()
	// 登录登出刷新token无需鉴权
	router.POST("/login", authMiddleware.LoginHandler)
	//router.GET("/login", authMiddleware.LoginHandler)
	router.POST("/logout", authMiddleware.LogoutHandler)
	router.POST("/refreshToken", authMiddleware.RefreshHandler)
	//}
	return router
}
