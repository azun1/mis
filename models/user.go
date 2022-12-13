package models

import (
	"MIS/pkg/logging"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Uuid       string     `json:"uuid" gorm:"type:char(36);not null;unique_index"`
	Email      string     `json:"email" gorm:"type:varchar(254);unique_index"`
	Name       string     `json:"username" gorm:"type:varchar(20);not null;unique_index"  validate:"required"`
	Password   string     `json:"password" gorm:"not null"  validate:"required"`
	RealName   string     `json:"realName" gorm:"type:varchar(20);not null;default:''"`
	Gender     string     `json:"gender"`                    // 性别
	Birth      *time.Time `json:"birth" gorm:"default:NULL"` // 出生日期
	Picture    string     `json:"picture"`
	Power      int        `json:"power" gorm:"default:1"`         // 用户类型
	LastActive *time.Time `json:"lastActive" gorm:"default:NULL"` // 用户上次活跃时间
	LogoutTime *time.Time `json:"logoutTime" gorm:"default:NULL"` // 用户上次登出时间
}

// FindUserByName 根据用户名称查找用户
func FindUserByName(name string) (user User, err error) {
	result := db.First(&user, "name = ?", name)
	err = result.Error
	if err != nil {
		logging.Error("FindUserByName name: %v error: %v", name, err)
	}
	return
}

// FindUserByEmail 根据邮箱查找用户
func FindUserByEmail(email string) (user User, err error) {
	result := db.First(&user, "email = ?", email)
	err = result.Error
	if err != nil {
		logging.Error("FindUserByEmail email: %v error: %v", email, err)
	}
	return
}

// IsUserNameExist 判断用户名称是否已存在
func IsUserNameExist(name string) bool {
	var user User
	result := db.First(&user, "name = ?", name)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}

// FindUserByUuid 根据用户uuid查找用户
func FindUserByUuid(uuid string) (user User, err error) {
	result := db.First(&user, "uuid = ?", uuid)
	err = result.Error
	if err != nil {
		logging.Error("FindUserByUuid uuid: %v error: %v", uuid, err)
	}
	return
}

// UpdateLastActiveTime 更新用户上次活跃时间
func (u *User) UpdateLastActiveTime() error {
	result := db.Model(u).Update("last_active", time.Now())
	return result.Error
}

// UpdateUserTime 更新用户上次登出时间
func (u *User) UpdateUserTime() error {
	result := db.Model(u).Update("logout_time", time.Now())
	return result.Error
}

// Add 添加新用户
func (u *User) Add() error {
	// todo 将用户头像上传至oss

	result := db.Create(u)
	if result.Error != nil {
		logging.Error("Common.Add user: %v error: %v", u, result.Error)
		return result.Error
	} else {
		logging.Info("Common.Add user: %v", u)
	}

	return nil
}

// Update 更新用户信息
func (u *User) Update(n User) error {
	result := db.Model(&User{}).
		Where("id = ?", u.ID).
		Select("Email", "Name", "RealName", "Gender", "Birth").
		Updates(n)
	if result.Error != nil {

	}
	db.Model(&User{}).
		Where("id = ?", u.ID).
		First(u)
	return result.Error
}

// Delete 删除用户
func (u *User) Delete() error {
	result := db.Delete(u)
	if result.Error != nil {
		logging.Error("Common.Delete user: %v error: %v", u, result.Error)
		return result.Error
	} else {
		logging.Info("Common.Delete user: %v", u)
	}

	return nil
}
