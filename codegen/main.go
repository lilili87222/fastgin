package codegen

import "codegen/config"

func main() {
	// 加载配置文件到全局配置结构体
	config.InitConfig()
	// 初始化数据库
	config.InitDatabase()
}
