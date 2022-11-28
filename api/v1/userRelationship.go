package v1

import (
	"MIS/api"
	"MIS/models"
	"MIS/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/mail"
)

type targetUuidForm struct {
	Uuid string `json:"uuid" validate:"required"`
}

type userRelationshipInfoForm struct {
	FamilyUuid       string `json:"family_uuid" validate:"required"`
	RelationshipType string `json:"relationship_type" validate:"required"`
	RelationshipInfo string `json:"relationship_info" validate:"required"`
}

// RequestConnect 发送(1条)关联账号请求
func RequestConnect(c *gin.Context) {
	// 基于输入的email / username 查找对应用户
	var user = api.CurrentUser(c)
	var userSelectForm struct {
		// 搜索依据
		Basis string `json:"basis" validate:"required"`
	}
	if !util.BindAndValid(c, &userSelectForm) {
		return
	}
	_, err := mail.ParseAddress(userSelectForm.Basis)
	if err != nil {
		// 用户输入的可能是username
		anotherUser, err := models.FindUserByName(userSelectForm.Basis)
		if err != nil {
			code := -1
			c.JSON(code, gin.H{
				"code":    code,
				"message": "该用户不存在",
			})
			return
		}
		err = user.RequestConnect(anotherUser.Uuid)
		if err != nil {
			api.ErrHandle(c, err)
			return
		}
	} else {
		// 用户输入的是邮箱
		// TODO: 需要User.go中添加FindUserByEmail
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "关联申请已发送",
	})
}

// AcceptConnect 同意(1条)关联账号请求
func AcceptConnect(c *gin.Context) {
	var user = api.CurrentUser(c)
	tarUuidForm := targetUuidForm{}
	if !util.BindAndValid(c, &tarUuidForm) {
		return
	}
	err := user.AcceptConnect(tarUuidForm.Uuid)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "已同意关联",
	})
}

// DeleteConnection 拒绝/删除(1条)关联关系
func DeleteConnection(c *gin.Context) {
	var user = api.CurrentUser(c)
	tarUuidForm := targetUuidForm{}
	if !util.BindAndValid(c, &tarUuidForm) {
		return
	}
	err := user.DeleteConnection(tarUuidForm.Uuid)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "已取消关联",
	})
}

// GetRelatedAccList 获取(同意/未同意的)关联账号列表
func GetRelatedAccList(c *gin.Context) {
	var user = api.CurrentUser(c)
	var uuidList = make([]string, 5)
	err := user.GetRelatedList(&uuidList)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":      http.StatusOK,
		"uuid_list": uuidList,
	})
}

// GetRelatedAccount 获取某个已关联账号的信息(关系类型, 备注)
func GetRelatedAccount(c *gin.Context) {
	var user = api.CurrentUser(c)
	tarUuidForm := targetUuidForm{}
	if !util.BindAndValid(c, &tarUuidForm) {
		return
	}
	var info = make([]string, 2)
	err := user.GetRelatedAccount(tarUuidForm.Uuid, &info)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":              http.StatusOK,
		"relationship_type": info[0],
		"relationship_info": info[1],
	})
}

// SetRelatedAccount 设置某个已关联账号的信息(关系类型, 备注)
func SetRelatedAccount(c *gin.Context) {
	var user = api.CurrentUser(c)
	relationshipInfoForm := userRelationshipInfoForm{}
	if !util.BindAndValid(c, &relationshipInfoForm) {
		return
	}
	var info = make([]string, 2)
	info[0] = relationshipInfoForm.RelationshipType
	info[1] = relationshipInfoForm.RelationshipInfo
	err := user.SetRelatedAccount(relationshipInfoForm.FamilyUuid, &info)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "已更新备注",
	})
}
