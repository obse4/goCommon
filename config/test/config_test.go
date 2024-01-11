package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/obse4/goCommon/config"
)

func TestConfig(t *testing.T) {
	type OtherConfig struct {
		IsTrue bool
		Type   string
	}
	var GlobalConfig struct {
		Name  string
		Id    int
		Other OtherConfig
		List  []OtherConfig
		Arr   []string
	}
	exePath, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exeDir := filepath.Dir(exePath)
	config.InitConfig(fmt.Sprintf("%s/config/env.yml", exeDir), &GlobalConfig)
	config.InitConfig("", &GlobalConfig)
}
