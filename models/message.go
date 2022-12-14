package models

import (
	"MIS/pkg/logging"
	"gorm.io/gorm"
	"time"
)

type MessageType int
type SenderType int

const (
	Text MessageType = iota
	Image
)
const (
	Common SenderType = iota // 此处原来的User与User表结构体名冲突，更改为Common
	Admin
	System
)

type MessageRecord struct {
	gorm.Model
	SenderUId string
	MessageType
	Content   string
	TargetUId string
}

type SystemNotice struct {
	gorm.Model
	SenderType
	Content string
	MessageType
	TargetUId string
}

// GetMessageList get message list through user_id and target_id
func GetMessageList(userId, targetId string) ([]MessageRecord, error) {
	var messageRecords []MessageRecord
	err := db.Order("created_at desc").Where("sender_uid = ? AND target_uid = ?", userId, targetId).Find(&messageRecords).Error
	if err != nil {
		logging.Info(err)
		return nil, err
	}
	return messageRecords, nil
}

// GetMessageByType returns all message divide by message type
func GetMessageByType(messageType int) ([]MessageRecord, error) {
	var messageRecords []MessageRecord
	err := db.Order("created_at desc").Where("message_type = ?", messageType).Find(&messageRecords).Error
	if err != nil {
		logging.Info(err)
		return nil, err
	}
	return messageRecords, nil
}

// GetMessageDetailByType first find all records divide by message type,then returns a string slice contains all messages.
func GetMessageDetailByType(messageType int) ([]string, error) {
	var messageRecords []MessageRecord
	err := db.Order("created_at desc").Where("message_type = ?", messageType).Find(&messageRecords).Error
	if err != nil {
		logging.Info(err)
		return nil, err
	}
	result := make([]string, 0)
	for i := 0; i < len(messageRecords); i++ {
		result = append(result, messageRecords[i].Content)
	}
	return result, nil
}

// DelMessageRecord delete message record by user_id target_id and create_time(created_at automatically generated by gorm)
func DelMessageRecord(userId, targetId string, createTime time.Time) error {
	err := db.Delete(&MessageRecord{}).Where("user_id = ? AND target_id = ? AND created_at = ?", userId, targetId, createTime).Error
	if err != nil {
		logging.Info(err)
		return err
	}
	return nil
}

// SaveMessageRecord save record into database.
func SaveMessageRecord(record MessageRecord) error {
	err := db.Create(&record).Error
	if err != nil {
		logging.Info(err)
		return err
	}
	return nil
}
