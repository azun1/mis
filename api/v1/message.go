package v1

import (
	"MIS/models"
	"MIS/pkg/e"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListResponse struct {
	Response
	MessageRecords []models.MessageRecord
}

type Response struct {
	StatusCode int32  `json:"status_code,omitempty"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

//type RecordDetailResponse struct {
//	Response
//	models.MedicalRecord
//}

func GetMessageList(c *gin.Context) {
	userId := c.Query("user_id")
	targetId := c.Query("target_id")
	result, err := models.GetMessageList(userId, targetId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: http.StatusInternalServerError,
			StatusMsg:  "Database Error!",
		})
	} else {
		c.JSON(http.StatusOK, ListResponse{
			Response{
				StatusCode: e.SUCCESS,
				StatusMsg:  "Get message records successfully",
			},
			result,
		})
	}

}

//type RecordResponse struct {
//	Response
//	ResponseList []models.MedicalRecord
//}
//
//func GetRecordList(c *gin.Context) {
//
//	result, err := models.GetMedicalRecordList()
//	if err != nil {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database error!"})
//		return
//	}
//
//	c.JSON(http.StatusOK, RecordResponse{
//		Response:     Response{StatusCode: 0},
//		ResponseList: result,
//	})
//}
//
//func GetDetail(c *gin.Context) {
//	id := c.Query("user_id")
//	uuid, err := strconv.ParseInt(id, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "uuid error!"})
//		return
//	}
//	result, err := models.GetMedicalRecordById(uuid)
//	if err != nil {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database error!"})
//		return
//	}
//	c.JSON(http.StatusOK, RecordResponse{
//		Response:     Response{StatusCode: 0},
//		ResponseList: result,
//	})
//}
//
//func DeleteRecord(c *gin.Context) {
//	id := c.Query("user_id")
//	uuid, err := strconv.ParseInt(id, 10, 64)
//	if err != nil {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "uuid error!"})
//		return
//	}
//	err = models.DeleteMedicalRecordById(uuid)
//	if err != nil {
//		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database error!"})
//		return
//	}
//	c.JSON(http.StatusOK, RecordResponse{
//		Response: Response{StatusCode: 0},
//	})
//}
