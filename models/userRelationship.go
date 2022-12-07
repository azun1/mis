package models

import (
	"MIS/pkg/logging"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UserRelationship 账户之间的多对多关系, 建立两个联合索引
type UserRelationship struct {
	gorm.Model
	// 联合索引遵循左前缀原则, priority 数值小的在左边(默认值为10), 相同则按在 struct 中的次序排
	SelfUuid         string `json:"self_uuid" gorm:"default:'';not null;uniqueIndex:s&f_uuid_index,priority:1;uniqueIndex:f&s_uuid_index,priority:2"`
	FamilyUuid       string `json:"family_uuid" gorm:"default:'';not null;uniqueIndex:s&f_uuid_index,priority:2;uniqueIndex:f&s_uuid_index,priority:1"`
	RelationshipType string `json:"relationship_type" gorm:"default:'';not null"`
	RelationshipInfo string `json:"relationship_info" gorm:"default:'';not null"`
	IsConfirmed      bool   `json:"is_confirmed" gorm:"type:tinyint(1);default:0;not null"`
}

// AddUserRelationship 添加用户时使用
type AddUserRelationship struct {
	SelfUuid         string `json:"self_uuid"`
	FamilyUuid       string `json:"family_uuid"`
	RelationshipType string `json:"relationship_type"`
	RelationshipInfo string `json:"relationship_info"`
}

// createNewRecord 创建一行新的账号关联申请
func createNewRecord(relationship *AddUserRelationship) error {
	urs := UserRelationship{
		Model:            gorm.Model{},
		SelfUuid:         relationship.SelfUuid,
		FamilyUuid:       relationship.FamilyUuid,
		RelationshipType: relationship.RelationshipType,
		RelationshipInfo: relationship.RelationshipInfo,
	}

	// 不使用Create()方法, 避免重复申请时出错
	// 重复申请时, 对方只能看到最新一次的信息
	// Ref: https://blog.csdn.net/weixin_44202304/article/details/120707068
	res := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "self_uuid"}, {Name: "family_uuid"}},
		UpdateAll: true,
	}).Create(&urs)
	// 如果是对已经有的记录进行的修改, 数据库中该行记录的id不会变, 但是id计数器会+1, 不清楚是不是bug

	// 插入出错, 返回错误
	if res.Error != nil {
		logging.Error("Common.Add user_relationship between : %v and %v error: %v", urs.SelfUuid, urs.FamilyUuid, res.Error)
		return res.Error
	} else {
		logging.Info("Common.Add user_relationship between : %v and %v", urs.SelfUuid, urs.FamilyUuid)
	}

	return nil
}

// RequestConnect 发送(1条)关联账号请求
func (u *User) RequestConnect(relationship *AddUserRelationship) error {
	var self_family_rs = make([]UserRelationship, 0, 5) // 一个账号最多与5个其他账号发生联系
	var family_self_rs = make([]UserRelationship, 0, 5) // 检查别人关联自己

	db.Where("self_uuid = ?", u.Uuid).Find(&self_family_rs)
	db.Where("family_uuid = ?", u.Uuid).Find(&family_self_rs)

	// 如果已经存在关联, 不再重复插入
	for i := range self_family_rs {
		if self_family_rs[i].FamilyUuid == relationship.FamilyUuid {
			if self_family_rs[i].IsConfirmed {
				logging.Info("The connection already existed")
				return nil
			} else {
				// 已存在关联申请记录, 但是尚未通过
				// 有可能是重复申请, 但是对备注之类的信息进行了修改
				return createNewRecord(relationship)
			}
		}
	}

	// 检查是否为双向申请
	for i := range family_self_rs {
		if family_self_rs[i].SelfUuid == relationship.FamilyUuid {
			// 是双向, 直接同意关联
			return u.AcceptConnect(relationship.FamilyUuid)
		}
	}

	// 如果不存在关联, 且未达到关联数量限制
	if len(self_family_rs) < 5 && len(family_self_rs) < 5 {
		return createNewRecord(relationship)
	} else {
		logging.Error("The number of self_family_rs of %v has achieve top limitation", u)
	}

	return nil
}

// AcceptConnect 同意(1条)关联账号的请求
func (u *User) AcceptConnect(targetUuid string) error {
	relationship := UserRelationship{}
	// 使关联请求记录通过
	relationship.SelfUuid, relationship.FamilyUuid = targetUuid, u.Uuid
	res := db.Where("self_uuid = ? AND family_uuid = ?",
		relationship.SelfUuid, relationship.FamilyUuid).Take(&relationship)
	if relationship.IsConfirmed == true {
		logging.Info("The connection has been accepted")
		return nil
	}
	// 更新字段的值
	relationship.IsConfirmed = true
	res = db.Save(&relationship)

	// 本账号添加关联记录
	relationship = UserRelationship{
		// 重新赋值的目的是重新获取id, 避免主键冲突
		SelfUuid:    u.Uuid,
		FamilyUuid:  targetUuid,
		IsConfirmed: true,
	}
	// 软删除之后使用.Create()方法似乎不能插入新记录, 也不能将软删除状态恢复, 故使用Upsert方法
	res = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "self_uuid"}, {Name: "family_uuid"}},
		UpdateAll: true,
	}).Create(&relationship)

	if res != nil {
		logging.Error(res.Error)
	}

	return res.Error
}

// DeleteConnection 拒绝/删除(1条)关联关系
func (u *User) DeleteConnection(targetUuid string) error {
	relationship := UserRelationship{}
	res := db.Where("self_uuid = ? AND family_uuid = ?",
		targetUuid, u.Uuid).Take(&relationship)
	if relationship.IsConfirmed == true {
		// 本账号到对方的关联记录也删除
		res = db.Where("self_uuid = ? AND family_uuid = ?",
			u.Uuid, targetUuid).Delete(&UserRelationship{})
		if res != nil {
			logging.Error("error: %v when delete the related user's record", res.Error)
		}
	} else {
		// 拒绝关联关系
		logging.Info("Refuse one connection")
	}
	// (软)删除记录
	res = db.Delete(&relationship)
	if res != nil {
		logging.Error(res.Error)
		return res.Error
	}
	logging.Info("Delete one connection record")

	return nil
}

// GetRelatedList 获取(同意/未同意的)关联账号列表
// 返回其他账号申请关联本账号的条目
// @return: 查询结果的Error
func (u *User) GetRelatedList(confirmed *[]UserRelationship, unconfirmed *[]string) error {
	var unconfirmedList = make([]UserRelationship, 0, 5) // 一个账号最多关联5个其他账号
	// 其他账号申请关联本账号的时候, 本账号是作为family_uuid存储在数据库的
	res := db.Where("family_uuid = ? AND is_confirmed = ?",
		u.Uuid, false).Find(&unconfirmedList)
	if len(unconfirmedList) == 0 {
		logging.Info("There is no connection record for now")
		return gorm.ErrRecordNotFound
	}
	for i := range unconfirmedList {
		*unconfirmed = append(*unconfirmed, unconfirmedList[i].SelfUuid)
	}

	res = db.Where("self_uuid = ?", u.Uuid).Find(confirmed)

	return res.Error
}

// GetRelatedAccount 获取某个已关联账号的信息(关系类型, 备注)
// @return: 参数string切片传递, info[0]: 关系类型, info[1]: 备注
func (u *User) GetRelatedAccount(targetUuid string, info *[]string) error {
	relationship := UserRelationship{}
	res := db.Where("self_uuid = ? AND family_uuid = ?",
		u.Uuid, targetUuid).Take(&relationship)
	(*info)[0] = relationship.RelationshipType
	(*info)[1] = relationship.RelationshipInfo
	return res.Error
}

// SetRelatedAccount 设置某个已关联账号的信息(关系类型, 备注)
func (u *User) SetRelatedAccount(targetUuid string, info *[]string) error {
	relationship := UserRelationship{}
	res := db.Where("self_uuid = ? AND family_uuid = ?",
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

// CheckRelated 检查u.Uuid和targetUuid之间是否存在已同意的关联关系
func (u *User) CheckRelated(targetUuid string) error {
	res := db.Where("self_uuid = ? AND family_uuid = ? AND is_confirmed = ?",
		u.Uuid, targetUuid, true).Take(&UserRelationship{})
	return res.Error
}
