package routers

import v1 "MIS/api/v1"

func MessageRouter() {
	messages := Apiv1.Group("/message")
	{
		messages.GET("/getList", v1.GetMessageList)
	}
}
