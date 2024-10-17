package route

import (
	"fastgin/config"
	"fastgin/database"
	"fastgin/modules/app/controller"
	"fastgin/modules/sys/dao"
	"fastgin/modules/sys/model"
	"github.com/gin-gonic/gin"
)

func InitDictionary(r *gin.RouterGroup) gin.IRoutes {
	insertApi()
	insertMenu()
	return registRoutes(r)
}
func registRoutes(r *gin.RouterGroup) gin.IRoutes {
	controller := controller.NewDictionaryController()
	router := r.Group("/dictionary")
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
func insertApi() {
	apiDao := dao.ApiDao{}
	var apis = []model.Api{
		{
			Method:   "POST",
			Path:     "/dictionary/index",
			Category: "dictionary",
			Desc:     "新增字典",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/dictionary/index/:id",
			Category: "dictionary",
			Desc:     "获取字典",
			Creator:  "系统",
		},
		{
			Method:   "PATCH",
			Path:     "/dictionary/index/:id",
			Category: "dictionary",
			Desc:     "更新字典",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/dictionary/index/:id",
			Category: "dictionary",
			Desc:     "删除字典",
			Creator:  "系统",
		},
		{
			Method:   "GET",
			Path:     "/dictionary/index",
			Category: "dictionary",
			Desc:     "搜索字典",
			Creator:  "系统",
		},
		{
			Method:   "DELETE",
			Path:     "/dictionary/index",
			Category: "dictionary",
			Desc:     "批量删除字典",
			Creator:  "系统",
		},
	}

	newRoleCasbin := make([]model.RoleCasbin, 0)
	for _, api := range apis {
		oldApi, _ := apiDao.GetApiDescByPath(api.Path, api.Method)
		if oldApi.Id == 0 {
			database.Create(&api)
			newRoleCasbin = append(newRoleCasbin, model.RoleCasbin{
				Keyword: "admin",
				Path:    api.Path,
				Method:  api.Method,
			})
		}
	}
	if len(newRoleCasbin) > 0 {
		rules := make([][]string, 0)
		for _, c := range newRoleCasbin {
			rules = append(rules, []string{
				c.Keyword, c.Path, c.Method,
			})
		}
		isAdd, err := config.CasbinEnforcer.AddPolicies(rules)
		if !isAdd {
			config.Log.Errorf("写入casbin数据失败：%v", err)
		}
	}

}
func insertMenu() {
	exist := false
	menuList, _ := database.ListAll[model.Menu]()
	for _, menu := range menuList {
		if menu.Component == "/app/dictionary/index" {
			exist = true
		}
	}
	if !exist {
		roles := []model.Role{
			{
				Model:   model.Model{Id: 1},
				Name:    "管理员",
				Keyword: "admin",
				Desc:    new(string),
				Sort:    1,
				Status:  1,
				Creator: "系统",
			},
		}
		appMenuId := uint(8)
		menus := []model.Menu{
			{
				Name:      "Dictionary",
				Title:     "字典管理",
				Icon:      nil,
				Path:      "dictionary",
				Component: "/app/dictionary/index",
				Sort:      11,
				ParentId:  &appMenuId,
				Roles:     roles,
				Creator:   "系统",
			},
		}
		database.Create(menus)
	}
}
