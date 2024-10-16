package generator

import (
	"fmt"
	"github.com/glebarez/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// 全局mysql数据库变量
var DB *gorm.DB

func InitDatabase() {
	if Instance.Database.Type == "mysql" {
		initMysql()
	} else if Instance.Database.Type == "sqlite" {
		initSqlLite()
	} else {
		panic(fmt.Errorf("mysql and sqllite support by default,不支持的数据库类型: %s", Instance.Database.Type))
	}
}

// 初始化mysql数据库
func initMysql() {
	mysqlConfig := Instance.Database.MysqlConfig
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
	sqlConfig := Instance.Database.SqlLiteConfig
	db, err := gorm.Open(sqlite.Open(sqlConfig.FilePath), &gorm.Config{
		// 禁用外键(指定外键时不会在sqlite创建真实的外键约束)
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {

		panic(fmt.Errorf("初始化sqlite数据库异常: %v", err))
	}
	DB = db
	createTables()

}

// 自动迁移表结构
func createTables() {

}
