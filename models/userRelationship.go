package models

import (
	"MIS/pkg/logging"
	"errors"
	"gorm.io/gorm"
)

// UserRelationship 账户之间的多对多关系, 建立两个联合索引
type UserRelationship struct {
	gorm.Model
	// 联合索引遵循左前缀原则, priority 数值小的在左边(默认值为10), 相同则按在 struct 中的次序排
	SelfUuid         string `json:"self_uuid" gorm:"not null;uniqueIndex:s&f_uuid_index,priority:1;uniqueIndex:f&s_uuid_index,priority:2"`
	FamilyUuid       string `json:"family_uuid" gorm:"not null;uniqueIndex:s&f_uuid_index,priority:2;uniqueIndex:f&s_uuid_index,priority:1"`
	RelationshipType string `json:"relationship_type" gorm:"default:'';not null"`
	RelationshipInfo string `json:"relationship_info" gorm:"default:'';not null"`
}

// createNewRecord 添加新的(账号关系)记录
func createNewRecord(selfUuid, targetUuid string) error {
	urs := UserRelationship{
		Model:      gorm.Model{},
		SelfUuid:   selfUuid,
		FamilyUuid: targetUuid,
	}
	res := db.Create(&urs)
	// 插入出错, 返回错误
	if res.Error != nil {
		logging.Error("Common.Add user_relationship between : %v and %v error: %v", selfUuid, targetUuid, res.Error)
		return res.Error
	} else {
		logging.Info("Common.Add user_relationship between : %v and %v", selfUuid, targetUuid)
	}

	return nil
}

// Connect 关联账号
func (u *User) Connect(targetUuid string) error {
	var relationships = make([]UserRelationship, 5) // 一个账号最多关联5个其他账号
	res := db.Where("selfUuid = ?", u.Uuid).Find(&relationships)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// 如果该账号还没有任何关联记录, 直接插入
		return createNewRecord(u.Uuid, targetUuid)
	}
	// 如果已经存在关联, 不再重复插入
	for i := range relationships {
		// 只获取下标, 避免结构体复制
		if relationships[i].FamilyUuid == targetUuid {
			logging.Info("The connection already existed")
			return nil
		}
	}
	// 如果未达到关联数量限制
	if len(relationships) < 5 {
		return createNewRecord(u.Uuid, targetUuid)
	} else {
		logging.Error("The number of relationships of %v has achieve top limitation", u)
	}

	return nil
}

// GetRelatedList 获取已关联账号列表
// @return: 参数返回uuid切片引用, 查询结果的Error
func (u *User) GetRelatedList(uuidList *[]string) error {
	var relationships = make([]UserRelationship, 5) // 一个账号最多关联5个其他账号
	res := db.Where("selfUuid = ?", u.Uuid).Find(&relationships)
	for i := range relationships {
		*uuidList = append(*uuidList, relationships[i].FamilyUuid)
	}
	return res.Error
}

// DisConnect 解除关联关系
func (u *User) DisConnect(targetUuid string) error {

	return nil
}

// SetRelatedAccount 设置某个已关联账号的信息(关系类型, 备注)
func (u *User) SetRelatedAccount(targetUuid string) error {

	return nil
}

// GetRelatedAccount 获取某个已关联账号的信息(关系类型, 备注)
// @return: 参数string切片传递, info[0]: 关系类型, info[1]: 备注
func (u *User) GetRelatedAccount(targetUuid string, info *[]string) error {

	return nil
}
