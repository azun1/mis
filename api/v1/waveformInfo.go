package v1

import (
	"MIS/api"
	"MIS/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AllMineHeartRate 获得自己所有的心率数据
func AllMineHeartRate(c *gin.Context) {
	var user = api.CurrentUser(c)
	var filePaths []string
	err := user.DownloadCSVByType("HeartRate", &filePaths)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	for i := range filePaths {
		// 使用Gin提供的文件下载服务
		c.File(filePaths[i])
	}
}

// AllMineBreathRate 获得自己的所有呼吸率数据
func AllMineBreathRate(c *gin.Context) {
	var user = api.CurrentUser(c)
	var filePaths []string
	err := user.DownloadCSVByType("BreathRate", &filePaths)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	for i := range filePaths {
		// 使用Gin提供的文件下载服务
		c.File(filePaths[i])
	}

}

// LatestHeartRate 获得最近的心率数据 (10s)
func LatestHeartRate(c *gin.Context) {
	var user = api.CurrentUser(c)
	var desc = models.Description{
		WaveformType: "HeartRate",
	}
	err := user.GetLatestRate(&desc)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, desc)
}

// LatestBreathRate 获得最近的呼吸率数据 (10s)
func LatestBreathRate(c *gin.Context) {
	var user = api.CurrentUser(c)
	var desc = models.Description{
		WaveformType: "BreathRate",
	}
	err := user.GetLatestRate(&desc)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, desc)
}
