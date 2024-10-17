package route

import (
	"fastgin/modules/sys/controller"
	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.RouterGroup) gin.IRoutes {
	menuController := controller.NewMenuController()
	router := r.Group("/menu")
	{
		router.GET("/tree", menuController.GetMenuTree)

		router.GET("/index", menuController.List)
		router.POST("/index", menuController.Create)
		router.PATCH("/index/:menuId", menuController.Update)
		router.DELETE("/index", menuController.BatchDeleteByIds)

		router.GET("/user/:userId", menuController.GetUserMenusByUserId)
		router.GET("/user_tree/:userId", menuController.GetUserMenuTreeByUserId)
	}
	return r
}
