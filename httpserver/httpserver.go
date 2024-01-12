package httpserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/obse4/goCommon/logger"
)

type HttpServerConfig struct {
	// 端口号 不填使用默认端口号8080
	Port   string `yaml:"port"`
	Mode   string `yaml:"mode" default:"debug" example:"debug/release/test"`
	Router *gin.Engine
}

func NewHttpServer(conf *HttpServerConfig) *gin.Engine {
	gin.SetMode(conf.Mode)
	router := gin.New()

	// 写入recover重新拉起
	router.Use(Recover)

	// 写入跨域中间件
	router.Use(Cors())

	// debug级别httpserver日志
	router.Use(DebugLogger())

	// 端口
	if conf.Port == "" {
		conf.Port = ":8080"
	}

	conf.Router = router

	return router
}

// Recover
// 确保该中间件处于最上层
// 可以防止程序挂掉
func Recover(ctx *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			// 打印错误
			logger.Error("HttpServer recover %v", r)
			// 打印错误堆栈信息
			logger.Error("HttpServer recover stack %s", string(debug.Stack()))

			recoverByte, _ := json.Marshal(r)

			ctx.JSON(http.StatusInternalServerError, gin.H{
				"code":    500,
				"message": string(recoverByte),
				"data":    nil,
			})
			//终止后续接口调用，不加的话recover到异常后，还会继续执行接口里后续代码
			ctx.Abort()
		}
	}()

	//加载完 defer recover，继续后续接口调用
	ctx.Next()
}

// 跨域访问中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func DebugLogger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		ctx.Next()
		endTime := time.Now()

		// 执行时间
		accessTime := endTime.Sub(startTime)

		logger.Debug("HttpServer | %3d | %13v | %15s | %s | %s |", ctx.Writer.Status(),
			accessTime,
			ctx.ClientIP(),
			ctx.Request.Method,
			ctx.Request.RequestURI)
	}
}

func (h *HttpServerConfig) Init() {
	port := h.Port
	if !strings.Contains(h.Port, ":") {
		port = fmt.Sprintf(":%s", h.Port)
	}
	srv := &http.Server{
		Addr:    port,
		Handler: h.Router,
	}
	logger.Info("--- Server Start ---")
	logger.Info("--- Server Listen:%s ---", h.Port)
	go func() {

		// 开启goroutine启动服务
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Server Fatal: %v", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器，为关闭服务器操作设置一个5秒的超时
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// kill -9 发送 syscall.SIGKILL 信号，但是不能被捕获，所以不需要添加它
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞在此，当接收到上述两种信号时才会往下执行
	<-quit
	logger.Info("--- Shutting Down Server ---")
	// 创建一个10秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// 10秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过10秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Forced To Shutdown: %v", err)
	}

	logger.Info("--- Server Exiting ---")
	defer logger.Info("--- Server Closed ---")
}
