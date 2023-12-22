package database

import (
	"fmt"

	"github.com/obse4/goCommon/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type MysqlConfig struct {
	Name     string   // 自定义名称
	Username string   // 用户名
	Password string   // 密码
	Database string   // 数据库
	Url      string   // url地址
	Port     string   // 端口
	Charset  string   // 编码方式 不填默认utf8
	Db       *gorm.DB // 数据库指针
}

func InitMysqlConnect(database *MysqlConfig) {
	if database.Charset == "" {
		database.Charset = "utf8"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local", database.Username, database.Password, database.Url, database.Port, database.Database, database.Charset)

	var err error
	database.Db, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		SkipInitializeWithVersion: false,
		DefaultStringSize:         255,
		DefaultDatetimePrecision:  nil,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		DontSupportForShareClause: false,
	}), &gorm.Config{
		Logger:                 gormLogger.Default.LogMode(gormLogger.Silent), // 忽略慢sql日志
		PrepareStmt:            false,                                         // 关闭预加载
		SkipDefaultTransaction: true,                                          // 关闭gorm事务模式
	})

	if err != nil {
		logger.Fatal("Mysql 数据库 %s 连接失败:%s", database.Name, err.Error())
	}

	logger.Info("Mysql 数据库 %s 连接成功", database.Name)
}
