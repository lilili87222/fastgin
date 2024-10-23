package route

import (
	"fastgin/modules/sys/controller"
	"fastgin/modules/sys/middleware"
	"github.com/gin-gonic/gin"
)

func InitBaseRoutes(router *gin.RouterGroup) gin.IRoutes {
	authMiddleware := middleware.GetJwtMiddleware()
	router.POST("/login", authMiddleware.LoginHandler)

	baseController := &controller.BaseController{}
	router.GET("/captcha", baseController.Captcha)
	router.POST("/register", baseController.Register)
	router.GET("/verifycode", baseController.SendVerifyCode)
	return router
}
