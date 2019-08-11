package server

import (
	"fmt"
	"gops/backend/glo"
	"time"

	loghandler "gops/backend/pkg/middleware/loghandler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	sysutil "gops/backend/server/route/sys"
)

// Serve Run
func Serve() {
	gin.DisableConsoleColor()
	engine := gin.Default()
	// 使用自定义日志句柄
	if glo.Config.GopsAPI.Log.Mode == "prod" {
		engine.Use(loghandler.Logger())
	}
	// engine.Use(loghandler.Logger())
	// 允许使用跨域请求  全局中间件
	engine.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-Token"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 因为这个路由，不需要登录操作，所以放到这里
	engine.POST("/InitTable", sysutil.InitTable)
	engine.POST("/UserLogin", sysutil.UserLogin)

	// 加载用户管理路由
	sysutil.InitRoute(engine)

	portVal := fmt.Sprintf("0.0.0.0:%d", glo.Config.GopsAPI.ServerPort)
	engine.Run(portVal)
}
