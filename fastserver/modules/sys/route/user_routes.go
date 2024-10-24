package route

import (
	"fastgin/modules/sys/controller"
	"github.com/gin-gonic/gin"
)

// 注册用户路由
func InitUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	userController := controller.NewUserController()
	router := r.Group("/user")
	{
		router.GET("/index", userController.GetUsers)
		router.POST("/index", userController.CreateUser)
		router.PATCH("/index/:userId", userController.Update)
		router.DELETE("/index", userController.BatchDeleteUserByIds)

		router.GET("/info", userController.GetUserInfo)
		router.PUT("/changePwd", userController.ChangePwd)
		router.POST("/logout", userController.Logout)
		//router.POST("/refreshToken", middleware.GetJwtMiddleware().RefreshHandler)
	}
	return r
}
