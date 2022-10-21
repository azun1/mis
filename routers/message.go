package routers

import v1 "MIS/api/v1"

func MessageRouter() {
	messages := Apiv1.Group("/message")
	{
		messages.GET("/getDetail", v1.GetDetail)
		messages.GET("/getList", v1.GetRecordList)
		messages.POST("/delete", v1.DeleteRecord)
	}
}
