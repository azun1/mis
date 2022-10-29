package v1

import (
	"MIS/models"
	"MIS/pkg/e"
	"MIS/pkg/logging"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type ListResponse struct {
	Response
	MessageRecords []models.MessageRecord
}

type MessageListResponse struct {
	Response
	MessageList []string
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
		return
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

func GetMessageByType(c *gin.Context) {
	messageType := c.Query("message_type")
	typeParam, err := strconv.Atoi(messageType)
	if err != nil {
		logging.Info(err)
		return
	}
	result, err := models.GetMessageByType(typeParam)
	if err != nil {
		logging.Info(err)
		c.JSON(http.StatusInternalServerError, ListResponse{
			Response:       Response{500, "Failed to get records."},
			MessageRecords: nil,
		})
		return
	}
	c.JSON(http.StatusOK, ListResponse{
		Response:       Response{200, "Get records by message type successfully."},
		MessageRecords: result,
	})
}

func GetMessageDetailByType(c *gin.Context) {
	messageType := c.Query("message_type")
	typeParam, err := strconv.Atoi(messageType)
	if err != nil {
		logging.Info(err)
	}
	result, err := models.GetMessageDetailByType(typeParam)
	if err != nil {
		logging.Info(err)
		c.JSON(http.StatusInternalServerError, ListResponse{
			Response:       Response{500, "Failed to get records."},
			MessageRecords: nil,
		})
		return
	}
	c.JSON(http.StatusOK, MessageListResponse{
		Response:    Response{200, "Get records by message type successfully."},
		MessageList: result,
	})
}

func DelMessageRecord(c *gin.Context) {
	userId := c.Query("user_id")
	targetId := c.Query("target_id")
	createTime := c.Query("created_at")
	createAt, err := time.Parse("2006-01-02 15:04:05", createTime)
	if err != nil {
		logging.Info(err)
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "failed to convert time string.",
		})
		return
	}
	err = models.DelMessageRecord(userId, targetId, createAt)
	if err != nil {
		logging.Info(err)
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 2,
			StatusMsg:  "failed to delete record.",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		StatusCode: 200,
		StatusMsg:  "Successfully delete.",
	})
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
