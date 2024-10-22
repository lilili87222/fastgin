package main

import (
	"fastgin/boost"
	"fastgin/boost/config"
	"fastgin/common/storage"
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
	config.InitConfig()

	// 初始化日志
	config.InitLogger()

	// 初始化数据库
	database.InitDatabaseConnection()
	// 初始化casbin策略管理器
	config.InitCasbinEnforcer(database.DB)

	// 初始化Validator数据校验
	config.InitValidate()

	storage.InitStorage(config.Configs.Storage)

	boost.StartWebService()

}
