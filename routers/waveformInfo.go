package routers

import (
	v1 "MIS/api/v1"
	"MIS/middleware"
)

func WaveformInfoRouter() {
	waveformInfo := Apiv1.Group("/waveform_info")
	{
		// 获得自己的所有心率数据
		waveformInfo.GET("/heart_rate", middleware.JWT(), v1.AllMineHeartRate)

		// 获得自己的所有呼吸率数据
		waveformInfo.GET("/breath_rate", middleware.JWT(), v1.AllMineBreathRate)

		// <展示用> 获得最近的心率数据 (10s)
		waveformInfo.GET("/latest_heart_rate", middleware.JWT(), v1.LatestHeartRate)

		// <展示用> 获得最近的呼吸率数据 (10s)
		waveformInfo.GET("/latest_breath_rate", middleware.JWT(), v1.LatestBreathRate)

		// <展示用> 获得关联账号最近的心率数据
		waveformInfo.POST("/related_latest_heart_rate", middleware.JWT(), v1.RelatedLatestHeartRate)

		// <展示用> 获得关联账号最近的呼吸率数据
		waveformInfo.POST("/related_latest_breath_rate", middleware.JWT(), v1.RelatedLatestBreathRate)
	}
}
