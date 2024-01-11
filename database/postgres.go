package database

import (
	"fmt"

	"github.com/obse4/goCommon/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgresConnect(database *PostgresConfig) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", database.Url, database.Username, database.Password, database.Database, database.Port)

	var err error
	database.Db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal("Postgres 数据库 %s 连接失败:%s", database.Name, err.Error())
	}

	logger.Info("Postgres 数据库 %s 连接成功", database.Name)

	return database.Db
}

type PostgresConfig struct {
	Name     string   // 自定义名称
	Username string   // 用户名
	Password string   // 密码
	Database string   // 数据库
	Url      string   // url地址
	Port     string   // 端口
	Db       *gorm.DB // 数据库指针
}
