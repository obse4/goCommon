# goCommon httpserver

## 使用
```
import "github.com/obse4/goCommon/httpserver"

func main() {
    // 配置httpserver配置文件
    var httpServer = httpserver.HttpServerConfig{
		Port: "8080",
		Mode: "test",
	}
    // 创建httpserver服务
    httpserver.NewHttpServer(&httpserverConfig)

    // 注册swagger
    // 配合github.com/swaggo/swag使用
    // ! 注意引入swag init 后的init函数
    httpserver.RegisterSwagger(httpServer.Router)

    // 注册路由
    httpServer.Router.GET("health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "ok",
			"data":    nil,
		})
	})

    // 初始化http服务监听
    httpserverConfig.Init()
}

```