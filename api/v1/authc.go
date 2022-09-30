package v1

import (
	"MIS/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Test(context *gin.Context) {
	logging.Info("日志打印测试")
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "服务启动成功",
	})
}
