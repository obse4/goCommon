package test

import (
	"fmt"
	"net/http"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/obse4/goCommon/httpserver"
)

func TestInitHttpServer(t *testing.T) {
	var httpserverConfig = httpserver.HttpServerConfig{
		Port: "8080",
		Mode: "debug",
	}
	httpserver.NewHttpServer(&httpserverConfig)

	httpserver.RegisterSwagger(httpserverConfig.Router)

	httpserverConfig.Router.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "ok",
			"data":    nil,
		})
	})

	// ip使用率限制
	limiter := httpserver.NewIPRateLimiter(10, 10)
	httpserverConfig.Router.Use(httpserver.IPRateLimitMiddleware(limiter))

	// 在另一个goroutine中，模拟发送Ctrl+C操作，即发送SIGINT信号
	go func() {
		time.AfterFunc(5*time.Second, func() {
			fmt.Println("Sending SIGINT signal...")

			process, err := os.FindProcess(os.Getpid())

			if err != nil {
				fmt.Printf("Failed to find process: %v", err)
			}

			err = process.Signal(syscall.SIGINT) // 向当前进程发送SIGINT信号
			if err != nil {
				fmt.Printf("Failed to send SIGINT signal: %v", err)
			}
		})
	}()

	httpserverConfig.Init()

}
