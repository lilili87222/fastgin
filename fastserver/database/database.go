package database

import (
	"fastgin/config"
	"fastgin/modules/sys/model"
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 全局mysql数据库变量
var DB *gorm.DB

func InitDatabase() {
	config.Log.Infof("选中的数据库类型" + config.Instance.Database.Type)
	if config.Instance.Database.Type == "mysql" {
		initMysql()
	} else if config.Instance.Database.Type == "sqlite" {
		initSqlLite()
	} else {
		panic(fmt.Errorf("mysql and sqllite support by default,不支持的数据库类型: %s", config.Instance.Database.Type))
	}
}

// 初始化mysql数据库
func initMysql() {
	mysqlConfig := config.Instance.Database.MysqlConfig
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
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		// 禁用外键(指定外键时不会在mysql创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
		//// 指定表前缀
		//NamingStrategy: schema.NamingStrategy{
		//	TablePrefix: config.mysqlConfig.TablePrefix + "_",
		//},
	})
	if err != nil {
		config.Log.Panicf("初始化mysql数据库异常: %v", err)
		panic(fmt.Errorf("初始化mysql数据库异常: %v", err))
	}
	// 开启mysql日志
	if mysqlConfig.LogMode {
		db.Debug()
	}
	DB = db
	createTables()

}

// 初始化sqllite数据库
func initSqlLite() {
	var sqlConfig = config.Instance.Database.SqlLiteConfig
	db, err := gorm.Open(sqlite.Open(sqlConfig.FilePath), &gorm.Config{
		// 禁用外键(指定外键时不会在sqlite创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		config.Log.Panicf("初始化sqlite数据库异常: %v", err)
		panic(fmt.Errorf("初始化sqlite数据库异常: %v", err))
	}
	DB = db
	createTables()
	//config.Log.Infof("初始化sqlite数据库完成!")
}

// 自动迁移表结构
func createTables() {
	if config.Instance.Database.CreateTables {
		config.Log.Infof("初始化数据库完成!")
		DB.AutoMigrate(
			&model.User{},
			&model.Role{},
			&model.Menu{},
			&model.Api{},
			&model.OperationLog{},
		)
	}
}
