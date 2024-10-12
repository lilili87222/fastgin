package config

import (
	"fastgin/internal/model/sys"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 全局mysql数据库变量
var DB *gorm.DB

func InitDatabase() {
	if Conf.Database.Type == "mysql" {
		initMysql()
	} else if Conf.Database.Type == "sqlite" {
		initSqlLite()
	} else {
		panic(fmt.Errorf("mysql and sqllite support by default,不支持的数据库类型: %s", Conf.Database.Type))
	}
}

// 初始化mysql数据库
func initMysql() {
	mysqlConfig := Conf.Database.MysqlConfig
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		mysqlConfig.Username,
		mysqlConfig.Password,
		mysqlConfig.Host,
		mysqlConfig.Port,
		mysqlConfig.Database,
		mysqlConfig.Charset,
		mysqlConfig.Collation,
		mysqlConfig.Query,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
		//// 指定表前缀
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix: config.mysqlConfig.TablePrefix + "_",
		//},
	})
	if err != nil {
		Log.Panicf("初始化mysql数据库异常: %v", err)
		panic(fmt.Errorf("初始化mysql数据库异常: %v", err))
	}
	// 开启mysql日志
	if mysqlConfig.LogMode {
		db.Debug()
	}
	DB = db
	initDatabaseTableAndData()
	Log.Infof("初始化mysql数据库完成!")
}

// 初始化sqllite数据库
func initSqlLite() {
	sqlConfig := Conf.Database.SqlLiteConfig
	db, err := gorm.Open(sqlite.Open(sqlConfig.FilePath), &gorm.Config{
		// 禁用外键(指定外键时不会在sqlite创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		Log.Panicf("初始化sqlite数据库异常: %v", err)
		panic(fmt.Errorf("初始化sqlite数据库异常: %v", err))
	}
	DB = db
	initDatabaseTableAndData()
	Log.Infof("初始化sqlite数据库完成!")
}

// 自动迁移表结构
func initDatabaseTableAndData() {
	if Conf.Database.CreateTables {
		DB.AutoMigrate(
			&sys.User{},
			&sys.Role{},
			&sys.Menu{},
			&sys.Api{},
			&sys.OperationLog{},
		)
	}
	if Conf.Database.InitData {
		InitData()
	}
}
