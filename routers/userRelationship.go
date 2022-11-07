package routers

import (
	v1 "MIS/api/v1"
	"MIS/middleware"
)

func UserRelationshipRouter() {
	userRelationship := Apiv1.Group("/user-relationship")
	{
		// 申请(1条)关联账号
		userRelationship.POST("/request-connect", middleware.JWT(), v1.RequestConnect)

		// 同意(1条)关联账号
		userRelationship.POST("/accept-connect", middleware.JWT(), v1.AcceptConnect)

		// 解除(1条)关联关系
		userRelationship.DELETE("/delete-connect", middleware.JWT(), v1.DeleteConnection)

		// 获取(同意/未同意的)关联账号列表
		userRelationship.GET("/get-related-account-list", middleware.JWT(), v1.GetRelatedAccList)

		// 获取某个已关联账号的信息(关系类型, 备注)
		userRelationship.GET("/get-related-account", middleware.JWT(), v1.GetRelatedAccount)

		// 设置某个已关联账号的信息(关系类型, 备注)
		userRelationship.PUT("/set-related-account", middleware.JWT(), v1.SetRelatedAccount)
	}
}
