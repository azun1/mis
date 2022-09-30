package routers

import v1 "MIS/api/v1"

func AuthcRouter() {
	authc := Apiv1.Group("/authc")
	{
		// 测试
		authc.GET("/test", v1.Test)
	}
}
