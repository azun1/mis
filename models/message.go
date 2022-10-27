package models

import (
	"MIS/pkg/logging"
	"gorm.io/gorm"
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
	SenderType
	TargetUId string
}

type SystemNotice struct {
	gorm.Model
	SenderType
	MessageType
	TargetUId string
}

func GetMessageList(userId, targetId string) ([]MessageRecord, error) {
	var messageRecords []MessageRecord
	err := db.Order("created_at desc").Where("sender_uid = ? AND target_uid = ?", userId, targetId).Find(&messageRecords).Error
	if err != nil {
		logging.Info(err)
		return messageRecords, err
	}
	return messageRecords, nil
}

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
