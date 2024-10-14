package config

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"os"
)

// 系统配置，对应yml
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便
const AppVersion = "v1.0.0"

// 全局配置变量
var Instance = new(config)

type config struct {
	System   *SystemConfig `mapstructure:"system" json:"system"`
	Logs     *LogsConfig   `mapstructure:"logs" json:"logs"`
	Database *struct {
		Type          string         `mapstructure:"type" json:"type"`
		MysqlConfig   *MysqlConfig   `mapstructure:"mysql" json:"mysql"`
		SqlLiteConfig *SqlLiteConfig `mapstructure:"sqlite" json:"sqlite"`
	} `mapstructure:"database" json:"database"`
}

type SystemConfig struct {
	Mode          string `mapstructure:"mode" json:"mode"`
	UrlPathPrefix string `mapstructure:"url-path-prefix" json:"urlPathPrefix"`
	Port          string `mapstructure:"port" json:"port"`
}

type LogsConfig struct {
	Level      zapcore.Level `mapstructure:"level" json:"level"`
	Path       string        `mapstructure:"path" json:"path"`
	MaxSize    int           `mapstructure:"max-size" json:"maxSize"`
	MaxBackups int           `mapstructure:"max-backups" json:"maxBackups"`
	MaxAge     int           `mapstructure:"max-age" json:"maxAge"`
	Compress   bool          `mapstructure:"compress" json:"compress"`
}

type MysqlConfig struct {
	Username    string `mapstructure:"username" json:"username"`
	Password    string `mapstructure:"password" json:"password"`
	Database    string `mapstructure:"database" json:"database"`
	Host        string `mapstructure:"host" json:"host"`
	Port        int    `mapstructure:"port" json:"port"`
	Query       string `mapstructure:"query" json:"query"`
	LogMode     bool   `mapstructure:"log-mode" json:"logMode"`
	TablePrefix string `mapstructure:"table-prefix" json:"tablePrefix"`
	Charset     string `mapstructure:"charset" json:"charset"`
	Collation   string `mapstructure:"collation" json:"collation"`
}
type SqlLiteConfig struct {
	FilePath string `mapstructure:"file-path" json:"filePath"`
}
type GeneratorConfig struct {
	//create-view: true # 是否生成视图vue
	//tables: ["users"] # 要生成的代码的表
	//module-name: "user" # 模块名称
	//out-dir: code_generated # 保存路径
	//run-sql: true # 是否执行sql
	OutDir     string   `mapstructure:"out-dir" json:"outDir"`
	Tables     []string `mapstructure:"tables" json:"tables"`
	ModuleName string   `mapstructure:"module-name" json:"moduleName"`
	RunSql     bool     `mapstructure:"run-sql" json:"runSql"`
	CreateView bool     `mapstructure:"create-view" json:"createView"`
}

// 设置读取配置信息
func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s \n", err))
	}

	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "./")
	// 读取配置信息
	err = viper.ReadInConfig()

	// 热更新配置
	viper.WatchConfig()
	if err != nil {
		panic(fmt.Errorf("读取配置文件失败:%s \n", err))
	}
	// 将读取的配置信息保存至全局变量Conf
	if err := viper.Unmarshal(Instance); err != nil {
		panic(fmt.Errorf("初始化配置文件失败:%s \n", err))
	}
}
