package config

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

// 系统配置，对应yml
// viper内置了mapstructure, yml文件用"-"区分单词, 转为驼峰方便
const AppVersion = "v1.0.0"

// 全局配置变量
var Instance = new(config)

type config struct {
	Generator *GeneratorConfig `mapstructure:"generator" json:"generator"`
	Database  *struct {
		Type          string         `mapstructure:"type" json:"type"`
		MysqlConfig   *MysqlConfig   `mapstructure:"mysql" json:"mysql"`
		SqlLiteConfig *SqlLiteConfig `mapstructure:"sqlite" json:"sqlite"`
	} `mapstructure:"database" json:"database"`
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
	OutDir      string   `mapstructure:"out-dir" json:"outDir"`
	Tables      []string `mapstructure:"tables" json:"tables"`
	ModuleName  string   `mapstructure:"module-name" json:"moduleName"`
	RunSql      bool     `mapstructure:"run-sql" json:"runSql"`
	CreateView  bool     `mapstructure:"create-view" json:"createView"`
	TablePrefix string   `mapstructure:"table-prefix" json:"tablePrefix"`
	Module      string   `mapstructure:"module" json:"module"`
}

// 设置读取配置信息
func InitConfig() {
	workDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("读取应用目录失败:%s \n", err))
	}

	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir)
	viper.AddConfigPath("..")
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
