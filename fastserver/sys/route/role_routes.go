package route

import (
	"fastgin/sys/controller"
	"github.com/gin-gonic/gin"
)

func InitRoleRoutes(r *gin.RouterGroup) gin.IRoutes {
	roleController := controller.NewRoleController()
	router := r.Group("/role")
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/index", roleController.GetRoles)
		router.POST("/index", roleController.CreateRole)
		router.PATCH("/index/:roleId", roleController.UpdateRoleById)
		router.DELETE("/index", roleController.BatchDeleteRoleByIds)

		router.GET("/menus/:roleId", roleController.GetRoleMenusById)
		router.PATCH("/menus/:roleId", roleController.UpdateRoleMenusById)

		router.GET("/apis/:roleId", roleController.GetRoleApisById)
		router.PATCH("/apis/:roleId", roleController.UpdateRoleApisById)

	}
	return r
}