package config

import (
	"os"
	"path/filepath"

	"github.com/obse4/goCommon/logger"
	"gopkg.in/yaml.v3"
)

// path 配置文件绝对路径 例如：/config/global.yml
// config 配置struct指针
// 可配置环境变量CONFIG_PATH来控制配置文件路径，权重最高 // 例如/config/global.yml CONFIG_PATH=/config/global.yml
// InitConfig函数配置的path权重其次
// path为空，默认使用当前执行文件所在目录下的config.yml文件

func InitConfig(path string, config any) {
	if os.Getenv("CONFIG_PATH") != "" {
		path = os.Getenv("CONFIG_PATH")
	} else if path == "" {
		path = getExeDir() + "/config.yml"
	}

	data, err := os.ReadFile(path)

	if err != nil {
		logger.Fatal("Config read error:%s", err.Error())
	}

	err = yaml.Unmarshal(data, config)

	if err != nil {
		logger.Fatal("Config unmarshal json error:%s", err.Error())
	}

	logger.Info("Config read %s %+v\n", path, config)
}

func getExeDir() string {
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exeDir := filepath.Dir(exePath)
	return exeDir
}
