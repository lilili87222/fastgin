package route

import (
"{{.Module}}/controller"
"fastgin/modules/sys/model"
"fastgin/modules/sys/service"
"github.com/gin-gonic/gin"
)

func Init{{.ModelName}}(r *gin.RouterGroup) gin.IRoutes {
	groupName := "{{.LowModelName}}"
	insertApiAndMenu(groupName)

	controller := controller.New{{.ModelName}}Controller()
	router := r.Group(groupName)
	{
		router.POST("/index", controller.Create)
		router.GET("/index/:id", controller.GetByID)
		router.PATCH("/index/:id", controller.Update)
		router.DELETE("/index/:id", controller.Delete)

		router.GET("/index", controller.List)
		router.DELETE("/index", controller.DeleteBatch)
	}
	return r
}
func insertApiAndMenu(groupName string) {
	tableDesc := "{{.TableComment}}"
	if tableDesc == "" {
		tableDesc = groupName
	}
	var apis = []model.Api{
		{
			Method:   "POST",
			Path:     "/" + groupName + "/index",
			Category: groupName,
			Des:     "新增" + tableDesc,
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/" + groupName + "/index/:id",
			Category: groupName,
			Des:     "获取" + tableDesc,
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/" + groupName + "/index/:id",
			Category: groupName,
			Des:     "更新" + tableDesc,
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/" + groupName + "/index/:id",
			Category: groupName,
			Des:     "删除" + tableDesc,
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/" + groupName + "/index",
			Category: groupName,
			Des:     "搜索" + tableDesc,
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/" + groupName + "/index",
			Category: groupName,
			Des:     "批量删除" + tableDesc,
			Creator:  "系统",
		},
	}
	service.NewApiService().InsertApisToAdmin(apis)

	{{ if .GenerateFront }}
	menu := model.Menu{
		Name:      "{{.ModelName}}",
		Title:     tableDesc + "管理",
		Icon:      "",
		Path:      groupName,
		Component: "/app/" + groupName + "/index",
		Sort:      11,
		Creator:   "系统",
	}
	service.NewMenuService().InsertAppMenuToAdmin(menu)
	{{ end }}
}
