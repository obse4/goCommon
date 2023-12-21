package logger

import (
	"testing"
)

func TestInitLogger(t *testing.T) {
	InitLogger(&LogConfig{
		LogOut: true,
	})
	Info("Hello %s", "world")
}

func TestLoggerLevel(t *testing.T) {
	InitLogger(&LogConfig{
		LogLevel: WarnLevel,
		LogOut:   true,
		StayDay:  1,
	})

	Debug("Hello %s", "world")
	Info("Hello %s", "world")
	Warn("Hello %s", "world")
	Error("Hello %s", "world")
}
