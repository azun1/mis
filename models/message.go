package models

import (
	"MIS/pkg/logging"
)

type GenderType string

type MedicalRecord struct {
	//Sex               uint8   `json:"sex"`
	Weight            float64 `json:"weight"`
	Height            float64 `json:"height"`
	Age               int     `json:"age"`
	UnderlyingDisease string  `json:"underlying_disease"`
}

func GetMedicalRecordById(id int64) (medicalRecord []MedicalRecord, err error) {
	err = db.Where("u_uuid = ?", id).Find(&medicalRecord).Error
	if err != nil {
		logging.Info("Unexpected error occurred when get medical records by field")
	}
	return
}

func GetMedicalRecordList() (medicalRecord []MedicalRecord, err error) {
	err = db.Find(&medicalRecord).Error
	if err != nil {
		logging.Info("Unexpected error occurred when get medical records")
	}
	return
}
func DeleteMedicalRecordById(id int64) (err error) {
	err = db.Where("u_uuid = ?", id).Delete(&MedicalRecord{}).Error
	if err != nil {
		logging.Info("Unexpected error occurred when delete medical records by uuid")
	}
	return
}
