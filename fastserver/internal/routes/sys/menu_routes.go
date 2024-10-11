package sys

import (
	"fastgin/internal/controller/sys"
	"github.com/gin-gonic/gin"
)

func InitMenuRoutes(r *gin.RouterGroup) gin.IRoutes {
	menuController := sys.NewMenuController()
	router := r.Group("/menu")
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/tree", menuController.GetMenuTree)
		router.GET("/list", menuController.GetMenus)
		router.POST("/create", menuController.CreateMenu)
		router.PATCH("/update/:menuId", menuController.UpdateMenuById)
		router.DELETE("/delete/batch", menuController.BatchDeleteMenuByIds)
		router.GET("/access/list/:userId", menuController.GetUserMenusByUserId)
		router.GET("/access/tree/:userId", menuController.GetUserMenuTreeByUserId)
	}

	return r
}
