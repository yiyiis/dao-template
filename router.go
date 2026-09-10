package main

import (
	"dao-template/api"
	. "dao-template/pkg/apiwarp"
	"github.com/gin-gonic/gin"
	"time"
)

func RegisterRouter(engine *gin.Engine) {
	// 健康检查
	engine.GET("/api/ping", TimeOut(time.Second*2), Controller(api.Ping))

	// 认证模块（公开接口）
	engine.POST("/api/auth/login", TimeOut(time.Second*5), Controller(api.UserLogin))

	// 用户模块（受 JWT 中间件保护）
	engine.GET("/api/user/info", TimeOut(time.Second*3), Controller(api.UserInfoDetail))
	engine.POST("/api/user/password", TimeOut(time.Second*3), Controller(api.UpdatePassword))
}
