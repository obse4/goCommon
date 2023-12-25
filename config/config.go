package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/obse4/goCommon/logger"
)

// folder 配置文件在执行文件下的相对路径 例如：global
// config 全局配置struct
// 配置环境变量CONFIG_MODE来控制配置文件名，不需要加后缀
// 例如/config/global.yml CONFIG_MODE=global
func InitConfig(folder string, config any) {
	fileName := os.Getenv("CONFIG_MODE")
	if fileName == "" {
		fileName = "env"
	}
	v := viper.New()
	v.SetConfigName(fileName)
	v.SetConfigType("yml")
	configPath := getExeDir()
	if folder != "" {
		configPath = filepath.Join(getExeDir(), folder)
	}

	v.AddConfigPath(configPath)

	if err := v.ReadInConfig(); err != nil {
		logger.Fatal("Config config read error:%s", err.Error())
	}
	if err := v.Unmarshal(&config); err != nil {
		logger.Fatal("Config unmarshal json error:%s", err.Error())
	}
	logger.Info("Config %+v\n", config)
}

func getExeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)
	return exeDir
}
