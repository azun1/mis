package routers

import (
	v1 "MIS/api/v1"
	"MIS/middleware"
)

func UserRouter() {
	user := Apiv1.Group("/user")
	{
		// 注册
		user.POST("/register", v1.Register)
		// 登录
		user.POST("/login", v1.Login)
		// 登出
		user.POST("/logout", middleware.JWT(), v1.Logout)
		// 注销
		user.DELETE("/delete", middleware.JWT(), v1.Delete)
		// 更新用户信息
		user.POST("/update", middleware.JWT(), v1.UpdateUserInfo)
	}
}
