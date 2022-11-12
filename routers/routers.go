package routers

import (
	"MIS/middleware"
	"MIS/pkg/settings"
	"github.com/gin-gonic/gin"
)

var Engine *gin.Engine

var Apiv1 *gin.RouterGroup

func RoutesController() *gin.Engine {
	Engine = gin.Default()
	Engine.Use(middleware.AllowAll())
	// 设置gin的工作模式
	gin.SetMode(settings.ServerSettings.RunMode)

	// v1版本路由
	Apiv1 = Engine.Group("api/v1")
	{
		MessageRouter()
		UserRouter()
	}

	return Engine
}
