package sys

import (
	"fastgin/internal/controller/sys"
	"github.com/gin-gonic/gin"
)

func InitRoleRoutes(r *gin.RouterGroup) gin.IRoutes {
	roleController := sys.NewRoleController()
	router := r.Group("/role")
	// 开启casbin鉴权中间件
	//router.Use(middleware.CasbinMiddleware())
	{
		router.GET("/list", roleController.GetRoles)
		router.POST("/create", roleController.CreateRole)
		router.PATCH("/update/:roleId", roleController.UpdateRoleById)
		router.GET("/menus/get/:roleId", roleController.GetRoleMenusById)
		router.PATCH("/menus/update/:roleId", roleController.UpdateRoleMenusById)
		router.GET("/apis/get/:roleId", roleController.GetRoleApisById)
		router.PATCH("/apis/update/:roleId", roleController.UpdateRoleApisById)
		router.DELETE("/delete/batch", roleController.BatchDeleteRoleByIds)
	}
	return r
}
