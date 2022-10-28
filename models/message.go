package models

import (
	"MIS/pkg/logging"
	"gorm.io/gorm"
	"time"
)

type MessageType int
type SenderType int

const (
	Image MessageType = iota
	Text
)
const (
	User SenderType = iota
	Admin
	System
)

type MessageRecord struct {
	gorm.Model
	SenderUId string
	MessageType
	Message string
	SenderType
	TargetUId string
}

type SystemNotice struct {
	gorm.Model
	SenderType
	Message string
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
		result = append(result, messageRecords[i].Message)
	}
	return result, nil
}

// DelMessageRecord delete message record by user_id target_id and create_time(created_at automaticly generated by gorm)
func DelMessageRecord(userId, targetId string, createTime time.Time) error {
	err := db.Delete(&MessageRecord{}).Where("user_id = ? AND target_id = ? AND created_at = ?", userId, targetId, createTime).Error
	if err != nil {
		logging.Info(err)
		return err
	}
	return nil
}

//func SaveMessageRecordForTest() {
//	messageRecord := MessageRecord{
//		Model:       gorm.Model{},
//		SenderUId:   "1",
//		MessageType: Text,
//		Message:     "hello this is a test message.",
//		SenderType:  User,
//		TargetUId:   "1",
//	}
//	db.Save(&messageRecord)
//}

//type GenderType string
//
//type MedicalRecord struct {
//	//Sex               GenderType   `json:"sex"`
//	Weight            float64 `json:"weight"`
//	Height            float64 `json:"height"`
//	Age               int     `json:"age"`
//	UnderlyingDisease string  `json:"underlying_disease"`
//}
//
//func GetMedicalRecordById(id int64) (medicalRecord []MedicalRecord, err error) {
//	err = db.Where("u_uuid = ?", id).Find(&medicalRecord).Error
//	if err != nil {
//		logging.Info("Unexpected error occurred when get medical records by field")
//	}
//	return
//}
//
//func GetMedicalRecordList() (medicalRecord []MedicalRecord, err error) {
//	err = db.Find(&medicalRecord).Error
//	if err != nil {
//		logging.Info("Unexpected error occurred when get medical records")
//	}
//	return
//}
//func DeleteMedicalRecordById(id int64) (err error) {
//	err = db.Where("u_uuid = ?", id).Delete(&MedicalRecord{}).Error
//	if err != nil {
//		logging.Info("Unexpected error occurred when delete medical records by uuid")
//	}
//	return
//}
