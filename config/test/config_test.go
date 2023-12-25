package test

import (
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
	}

	config.InitConfig("config", &GlobalConfig)
}
