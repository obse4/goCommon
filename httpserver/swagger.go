package httpserver

import (
	// _ "github.com/obse4/goCommon/httpserver/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 注意引入docs下的init函数
// 配合github.com/swaggo/swag使用
// go install github.com/swaggo/swag/cmd/swag@latest
// 生成swagger默认文件夹及根据注释写入内容 swag init
// 访问 http://127.0.0.1:8080/swagger/index.html#/
// ! 注意引入swag init 后的init函数
// import _ "your_project/docs"
func RegisterSwagger(r *gin.Engine) *gin.Engine {
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, func(c *ginSwagger.Config) { c.PersistAuthorization = true }))

	return r
}
