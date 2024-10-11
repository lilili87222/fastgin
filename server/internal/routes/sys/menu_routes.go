package sys

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go-web-mini/internal/controller/sys"
	"go-web-mini/internal/middleware"
)

func InitMenuRoutes(r *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) gin.IRoutes {
	menuController := sys.NewMenuController()
	router := r.Group("/menu")
	// 开启jwt认证中间件
	router.Use(authMiddleware.MiddlewareFunc())
	// 开启casbin鉴权中间件
	router.Use(middleware.CasbinMiddleware())
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
