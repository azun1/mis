package routers

import (
	v1 "MIS/api/v1"
)

// @Question: 这些接口需要做JWT鉴权吗
func UserRelationshipRouter() {
	userRelationship := Apiv1.Group("/user-relationship")
	{
		// 申请(1条)关联账号
		userRelationship.POST("/request-connect", v1.RequestConnect)

		// 同意(1条)关联账号
		userRelationship.POST("/accept-connect", v1.AcceptConnect)

		// 解除(1条)关联关系
		userRelationship.DELETE("/delete-connect", v1.DeleteConnection)

		// 获取(同意/未同意的)关联账号列表
		userRelationship.GET("/get-related-account-list", v1.GetRelatedAccList)

		// 获取某个已关联账号的信息(关系类型, 备注)
		userRelationship.GET("/get-related-account", v1.GetRelatedAccount)

		// 设置某个已关联账号的信息(关系类型, 备注)
		userRelationship.PUT("/set-related-account", v1.SetRelatedAccount)
	}
}
