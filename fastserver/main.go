package main

import (
	"fastgin/boost"
	config2 "fastgin/boost/config"
	"fastgin/database"
)

// @title Go Web fastgin API
// @version 1.0
// @description This is a sample server for a Go web mini project.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 192.168.123.214:8088
// @BasePath /
func main() {

	// 加载配置文件到全局配置结构体
	config2.InitConfig()

	// 初始化日志
	config2.InitLogger()

	// 初始化数据库
	database.InitDatabaseConnection()
	// 初始化casbin策略管理器
	config2.InitCasbinEnforcer(database.DB)

	// 初始化Validator数据校验
	config2.InitValidate()

	boost.StartWebService()

}
