package api

import (
	"MIS/models"
	"MIS/pkg/logging"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// ErrHandle 接口层笼统错误处理
func ErrHandle(context *gin.Context, err error) {
	logging.Error(err.Error())
	if err == gorm.ErrRecordNotFound {
		context.JSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	context.JSON(http.StatusInternalServerError, gin.H{
		"code": http.StatusInternalServerError,
		"err":  err.Error(),
	})
}

// CurrentUser 获取当前用户信息
func CurrentUser(context *gin.Context) *models.User {
	return context.MustGet("user").(*models.User)
}
