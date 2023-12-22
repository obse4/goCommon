package test

import (
	"os"
	"testing"
	"time"

	"github.com/obse4/goCommon/httpserver"
)

func TestInitHttpServer(t *testing.T) {
	var httpserverConfig = httpserver.HttpServerConfig{
		Port: "8080",
		Mode: "test",
	}
	httpserver.NewHttpServer(&httpserverConfig)

	httpserver.RegisterSwagger(httpserverConfig.Router)

	go func() {
		time.AfterFunc(5*time.Second, func() {
			os.Exit(0)
		})
	}()

	httpserverConfig.Init()

}
