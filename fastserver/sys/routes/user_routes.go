package routes

import (
	"fastgin/sys/controller"
	"fastgin/sys/middleware"
	"github.com/gin-gonic/gin"
)

// 注册用户路由
func InitUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	userController := controller.NewUserController()
	router := r.Group("/user")
	{
		router.GET("/index", userController.GetUsers)
		router.POST("/index", userController.CreateUser)
		router.PATCH("/index/:userId", userController.UpdateUserById)
		router.DELETE("/index", userController.BatchDeleteUserByIds)

		router.GET("/info", userController.GetUserInfo)
		router.PUT("/changePwd", userController.ChangePwd)
		router.POST("/logout", middleware.GetJwtMiddleware().LogoutHandler)
		router.POST("/refreshToken", middleware.GetJwtMiddleware().RefreshHandler)
	}
	return r
}
