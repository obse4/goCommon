package database

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/obse4/goCommon/logger"
)

type RedisConfig struct {
	Name     string      `yaml:"name"`     // 自定义名称
	Url      string      `yaml:"url"`      // url连接
	Port     string      `yaml:"port"`     // 端口
	Password string      `yaml:"password"` // 密码 非必填
	Pool     *redis.Pool // redis连接池
}

func InitRedisPool(database *RedisConfig) *redis.Pool {
	database.Pool = &redis.Pool{
		MaxIdle: 0,
		Wait:    true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%s", database.Url, database.Port), redis.DialPassword(database.Password))

			if err != nil {
				logger.Error("Redis error is %s", err.Error())
				return nil, err
			}
			return c, nil
		},
	}
	// 检查连接

	pong, err := redis.String(database.Pool.Get().Do("PING"))

	if err != nil {
		logger.Fatal("Redis %s 连接错误:%s", database.Name, err.Error())
	}

	if pong == "PONG" {
		logger.Info("Redis %s 连接成功", database.Name)
	} else {
		logger.Fatal("Redis %s 连接失败", database.Name)
	}

	return database.Pool
}

func GetRedisConn(pool *redis.Pool, db int) redis.Conn {
	conn := pool.Get()
	conn.Do("SELECT", db)
	return conn
}
