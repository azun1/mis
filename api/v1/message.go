package v1

import (
	"MIS/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type RecordDetailResponse struct {
	Response
	models.MedicalRecord
}

type RecordResponse struct {
	Response
	ResponseList []models.MedicalRecord
}

func GetRecordList(c *gin.Context) {

	result, err := models.GetMedicalRecordList()
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database error!"})
		return
	}

	c.JSON(http.StatusOK, RecordResponse{
		Response:     Response{StatusCode: 0},
		ResponseList: result,
	})
}

func GetDetail(c *gin.Context) {
	id := c.Query("user_id")
	uuid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "uuid error!"})
		return
	}
	result, err := models.GetMedicalRecordById(uuid)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database error!"})
		return
	}
	c.JSON(http.StatusOK, RecordResponse{
		Response:     Response{StatusCode: 0},
		ResponseList: result,
	})
}

func DeleteRecord(c *gin.Context) {
	id := c.Query("user_id")
	uuid, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "uuid error!"})
		return
	}
	err = models.DeleteMedicalRecordById(uuid)
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "database error!"})
		return
	}
	c.JSON(http.StatusOK, RecordResponse{
		Response: Response{StatusCode: 0},
	})
}
