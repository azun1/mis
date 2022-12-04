package v1

import (
	"MIS/models"
	"MIS/pkg/e"
	"MIS/pkg/logging"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

func GetMessageList(c *gin.Context) {
	var MessageListForm struct {
		UserId   string `json:"user_id"`
		TargetId string `json:"target_id"`
	}
	err := c.ShouldBind(&MessageListForm)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: http.StatusInternalServerError,
			StatusMsg:  "invalid request param!",
		})
		return
	}
	result, err := models.GetMessageList(MessageListForm.UserId, MessageListForm.TargetId)
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
	var TypeForm struct {
		MessageType string `json:"message_type"`
	}
	err := c.ShouldBindJSON(&TypeForm)
	typeParam, err := strconv.Atoi(TypeForm.MessageType)
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
	var TypeForm struct {
		Type string `json:"message_type"`
	}
	err := c.ShouldBindJSON(&TypeForm)
	typeParam, err := strconv.Atoi(TypeForm.Type)
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
		Response:    Response{200, "Get MessageDetail by message type successfully."},
		MessageList: result,
	})
}

func DelMessageRecord(c *gin.Context) {
	var DelMessageForm struct {
		UserId     string `json:"user_id" validate:"required"`
		TargetId   string `json:"target_id" validate:"required"`
		CreateTime string `json:"created_at" validate:"required"`
	}
	err := c.ShouldBindJSON(&DelMessageForm)
	fmt.Printf("%+v", DelMessageForm)
	createAt, err := time.Parse("2006-01-02 15:04:05", DelMessageForm.CreateTime)
	if err != nil {
		logging.Info(err)
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 1,
			StatusMsg:  "failed to convert time string.",
		})
		return
	}
	err = models.DelMessageRecord(DelMessageForm.UserId, DelMessageForm.TargetId, createAt)
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

func SaveRecord(c *gin.Context) {
	var SaveForm struct {
		UserId      string `json:"user_id"`
		TargetId    string `json:"target_id"`
		Content     string `json:"content"`
		MessageType string `json:"message_type"`
	}
	err := c.ShouldBindJSON(&SaveForm)
	if err != nil {
		logging.Info(err)
		return
	}
	Type, _ := strconv.Atoi(SaveForm.MessageType)
	var RecordType models.MessageType
	switch Type {
	case 1:
		RecordType = models.Text
	case 2:
		RecordType = models.Image
	default:
		c.JSON(http.StatusBadRequest, Response{
			StatusCode: 400,
			StatusMsg:  "Param error ",
		})
		return
	}
	messageRecord := models.MessageRecord{
		Model:       gorm.Model{},
		SenderUId:   SaveForm.UserId,
		MessageType: RecordType,
		Content:     SaveForm.Content,
		TargetUId:   SaveForm.TargetId,
	}
	err = models.SaveMessageRecord(messageRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			StatusCode: 500,
			StatusMsg:  "Failed to save record!",
		})
	} else {
		c.JSON(http.StatusOK, Response{
			StatusCode: 200,
			StatusMsg:  "Successfully save record!",
		})
	}
}
