package logger

import (
	"testing"
)

func TestInitLogger(t *testing.T) {
	InitLogger(&LogConfig{
		LogLevel: 0,
		LogOut:   true,
		// LogFile:  `C:\Users\tankk\Desktop\code\github.obse4\goCommon\log`,
	})
	Info("HELLO %s", "world")
}
