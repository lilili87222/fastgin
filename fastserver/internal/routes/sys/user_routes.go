package sys

import (
	"fastgin/internal/controller/sys"
	"fastgin/internal/middleware"
	"github.com/gin-gonic/gin"
)

// 注册用户路由
func InitUserRoutes(r *gin.RouterGroup) gin.IRoutes {
	userController := sys.NewUserController()
	router := r.Group("/user")
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/info", userController.GetUserInfo)
		//router.POST("/info", userController.GetUserInfo)
		router.GET("/list", userController.GetUsers)
		router.PUT("/changePwd", userController.ChangePwd)
		router.POST("/create", userController.CreateUser)
		router.PATCH("/update/:userId", userController.UpdateUserById)
		router.DELETE("/delete/batch", userController.BatchDeleteUserByIds)

		router.POST("/logout", middleware.GetJwtMiddleware().LogoutHandler)

		router.POST("/refreshToken", middleware.GetJwtMiddleware().RefreshHandler)
	}
	return r
}
