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

// AddUserForm 添加用户时使用
type AddUserForm struct {
	// 搜索依据
	Basis string `json:"basis" validate:"required"`
	// 前端可以显示为 关系类型 或 关系标签
	RelationshipType string `json:"relationship_type"`
	// 前端可以显示为 备注
	RelationshipInfo string `json:"relationship_info"`
}

// 描述关联关系的信息
type userRelationshipInfoForm struct {
	FamilyUuid       string `json:"family_uuid" validate:"required"`
	RelationshipType string `json:"relationship_type" validate:"required"`
	RelationshipInfo string `json:"relationship_info" validate:"required"`
}

// 描述存在关联关系的用户的信息
type userRelatedForm struct {
	Uuid             string `json:"uuid"`
	Name             string `json:"name"`
	Gender           string `json:"gender"`
	Email            string `json:"email"`
	RelationshipType string `json:"relationship_type"`
	RelationshipInfo string `json:"relationship_info"`
}

// RequestConnect 发送(1条)关联账号请求
func RequestConnect(c *gin.Context) {
	// 基于输入的email / username 查找对应用户
	var user = api.CurrentUser(c)
	var anotherUser models.User
	addUserForm := AddUserForm{}
	if !util.BindAndValid(c, &addUserForm) {
		return
	}
	_, err := mail.ParseAddress(addUserForm.Basis)
	if err != nil {
		// 用户输入的可能是username
		anotherUser, err = models.FindUserByName(addUserForm.Basis)
		if err != nil {
			code := -1
			c.JSON(code, gin.H{
				"code":    code,
				"message": "该用户不存在",
			})
			return
		}
	} else {
		// 用户输入的是邮箱
		anotherUser, err = models.FindUserByEmail(addUserForm.Basis)
		if err != nil {
			code := -1
			c.JSON(code, gin.H{
				"code":    code,
				"message": "该用户不存在",
			})
			return
		}
	}
	if anotherUser.Uuid == user.Uuid {
		c.JSON(http.StatusForbidden, gin.H{
			"message": "不能自相关联，请尝试关联到其他账号",
		})
		return
	}
	addUserRelationship := models.AddUserRelationship{
		SelfUuid:         user.Uuid,
		FamilyUuid:       anotherUser.Uuid,
		RelationshipType: addUserForm.RelationshipType,
		RelationshipInfo: addUserForm.RelationshipInfo,
	}
	err = user.RequestConnect(&addUserRelationship)
	if err != nil {
		api.ErrHandle(c, err)
		return
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
// 未关联的账号返回uuid, name, gender
// 已关联的账号返回uuid, name, gender, email, RelationshipType, RelationshipInfo
func GetRelatedAccList(c *gin.Context) {
	var user = api.CurrentUser(c)
	var confirmedRs = make([]models.UserRelationship, 0, 5)
	var unconfirmedUuidList = make([]string, 0, 5)

	// 参数返回已关联的账号(confirmed): {
	// 			FamilyUuid,
	//          RelationshipType,
	//          RelationshipInfo,
	//        }
	//       未同意的关联申请(unconfirmed): {
	//          FamilyUuid,
	//       }
	err := user.GetRelatedList(&confirmedRs, &unconfirmedUuidList)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	var confirmedList = make([]userRelatedForm, len(confirmedRs))
	for i := range confirmedList {
		cur, _ := models.FindUserByUuid(confirmedRs[i].FamilyUuid)
		confirmedList[i] = userRelatedForm{
			Uuid:             cur.Uuid,
			Name:             cur.Name,
			Gender:           cur.Gender,
			Email:            cur.Email,
			RelationshipType: confirmedRs[i].RelationshipType,
			RelationshipInfo: confirmedRs[i].RelationshipInfo,
		}
	}
	var unconfirmedList = make([]userRelatedForm, len(unconfirmedUuidList))
	for i := range unconfirmedList {
		cur, _ := models.FindUserByUuid(unconfirmedUuidList[i])
		unconfirmedList[i] = userRelatedForm{
			Uuid:   cur.Uuid,
			Name:   cur.Name,
			Gender: cur.Gender,
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":             http.StatusOK,
		"confirmed_list":   confirmedList,
		"unconfirmed_list": unconfirmedList,
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
