package test

import (
	"context"
	"testing"
	"time"

	"github.com/obse4/goCommon/logger"
)

func TestInitLogger(t *testing.T) {
	logger.InitLogger(&logger.LogConfig{
		LogOut: true,
	})
	logger.Info("Hello %s", "world")
}

func TestLoggerLevel(t *testing.T) {
	logger.InitLogger(&logger.LogConfig{
		LogLevel: logger.WarnLevel,
		LogOut:   true,
		StayDay:  1,
	})

	logger.Debug("Hello %s", "world")
	logger.Info("Hello %s", "world")
	logger.Warn("Hello %s", "world")
	logger.Error("Hello %s", "world")
}

func TestRunTime(t *testing.T) {
	logger.InitLogger(&logger.LogConfig{
		LogLevel: logger.DebugLevel,
		LogOut:   true,
		StayDay:  1,
	})

	ctx := logger.Time(context.TODO(), "TestRunTime")
	time.Sleep(time.Second)
	logger.TimeEnd(ctx, "TestRunTime").Debug()
}
