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
	SelfUuid         string `json:"self_uuid" gorm:"default:'';not null;uniqueIndex:s&f_uuid_index,priority:1;uniqueIndex:f&s_uuid_index,priority:2"`
	FamilyUuid       string `json:"family_uuid" gorm:"default:'';not null;uniqueIndex:s&f_uuid_index,priority:2;uniqueIndex:f&s_uuid_index,priority:1"`
	RelationshipType string `json:"relationship_type" gorm:"default:'';not null"`
	RelationshipInfo string `json:"relationship_info" gorm:"default:'';not null"`
	IsConfirmed      bool   `json:"is_confirmed" gorm:"type:tinyint(1);default:0"`
}

// createNewRecord 创建一行新的账号关联申请
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

// RequestConnect 发送(1条)关联账号请求
func (u *User) RequestConnect(targetUuid string) error {
	var relationships = make([]UserRelationship, 5) // 一个账号最多关联5个其他账号
	res := db.Where("selfUuid = ?", u.Uuid).Find(&relationships)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		// 如果该账号还没有任何关联记录, 直接插入
		// 插入的时候, 在被申请人那里留一条待同意记录, 申请方不添加记录
		return createNewRecord(targetUuid, u.Uuid)
	}
	// 如果已经存在关联, 不再重复插入
	for i := range relationships {
		// 只获取下标, 避免结构体复制
		if relationships[i].FamilyUuid == targetUuid {
			if relationships[i].IsConfirmed == false {
				relationships[i].IsConfirmed = true
				// 更新字段的值
				db.Save(&relationships[i])
				logging.Info("Confirmed the connection")
				return nil
			}
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

// AcceptConnect 同意(1条)关联账号的请求
func (u *User) AcceptConnect(targetUuid string) error {
	relationship := UserRelationship{}
	res := db.Where("selfUuid = ? AND Family_uuid = ?",
		u.Uuid, targetUuid).Take(&relationship)
	if relationship.IsConfirmed == true {
		logging.Info("The connection has been accepted")
		return nil
	}
	// 更新字段的值
	relationship.IsConfirmed = true
	res = db.Save(&relationship)
	// 给申请方添加记录
	relationship.SelfUuid, relationship.FamilyUuid = targetUuid, u.Uuid
	db.Create(&relationship)
	if res != nil {
		logging.Error(res.Error)
	}
	return res.Error
}

// DeleteConnection 拒绝/删除(1条)关联关系
func (u *User) DeleteConnection(targetUuid string) error {
	relationship := UserRelationship{}
	res := db.Where("selfUuid = ? AND Family_uuid = ?", u.Uuid, targetUuid).
		Take(&relationship)
	if relationship.IsConfirmed == true {
		// 被删除方相关记录也删除
		res = db.Where("selfUuid = ? AND Family_uuid = ?", targetUuid, u.Uuid).
			Delete(&UserRelationship{})
		if res != nil {
			logging.Error("error: %v when delete the related user's record", res.Error)
		}
	} else {
		// 拒绝关联关系
		logging.Info("Refuse one connection")
	}
	// 删除记录
	res = db.Delete(&relationship)
	if res != nil {
		logging.Error(res.Error)
		return res.Error
	}
	logging.Info("Delete one connection record")
	return nil
}

// GetRelatedList 获取(同意/未同意的)关联账号列表
// @return: 参数返回uuid切片引用, 查询结果的Error
func (u *User) GetRelatedList(uuidList *[]string) error {
	var relationships = make([]UserRelationship, 5) // 一个账号最多关联5个其他账号
	res := db.Where("selfUuid = ?", u.Uuid).Find(&relationships)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		logging.Info("There is no connection record for now")
		return gorm.ErrRecordNotFound
	}
	for i := range relationships {
		*uuidList = append(*uuidList, relationships[i].FamilyUuid)
	}
	return res.Error
}

// GetRelatedAccount 获取某个已关联账号的信息(关系类型, 备注)
// @return: 参数string切片传递, info[0]: 关系类型, info[1]: 备注
func (u *User) GetRelatedAccount(targetUuid string, info *[]string) error {
	relationship := UserRelationship{}
	res := db.Where("selfUuid = ? AND Family_uuid = ?",
		u.Uuid, targetUuid).Take(&relationship)
	(*info)[0] = relationship.RelationshipType
	(*info)[1] = relationship.RelationshipInfo
	return res.Error
}

// SetRelatedAccount 设置某个已关联账号的信息(关系类型, 备注)
func (u *User) SetRelatedAccount(targetUuid string, info *[]string) error {
	relationship := UserRelationship{}
	res := db.Where("selfUuid = ? AND Family_uuid = ?",
		u.Uuid, targetUuid).Take(&relationship)
	// 更新字段的值
	relationship.RelationshipType = (*info)[0]
	relationship.RelationshipInfo = (*info)[1]
	res = db.Save(&relationship)
	if res != nil {
		logging.Error(res.Error)
	}
	return res.Error
}
