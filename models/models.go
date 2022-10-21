package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	UUuid     uuid.UUID `json:"u_uuid"`
	UEmail    string    `gorm:"type: varchar(254)" json:"u_email"`
	UName     string    `gorm:"type: varchar(20)" json:"u_name"`
	UPassword string    `gorm:"type: varchar(128)" json:"u_password"`
	URealname string    `gorm:"type: varchar(20)" json:"u_realname"`
	UPicture  string    `gorm:"type: varchar(500)" json:"u_picture"`
}

type UserRelation struct {
	gorm.Model
	SelfUuid   uuid.UUID `json:"self_uuid"`
	FamilyUuid uuid.UUID `json:"family_uuid"`
}

type MedicalRecord struct {
	gorm.Model
	UUuid             uuid.UUID `json:"u_uuid"`
	Sex               string    `gorm:"type: enum('保密', '男', '女')" json:"sex"`
	Weight            float32   `gorm:"type: float" json:"weight"`
	Height            float32   `gorm:"type: float" json:"height"`
	Age               uint      `gorm:"type: tinyint unsigned" json:"age"`
	UnderlyingDisease string    `gorm:"type: varchar(1500)" json:"underlying_disease"`
}

type WaveformInfo struct {
	gorm.Model
	UUuid     uuid.UUID `json:"u_uuid"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	WaveType  string    `json:"wave_type"`
	FilePath  string    `json:"file_path"`
}
