package models

import (
	"MIS/pkg/logging"
	"MIS/pkg/settings"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// 全局数据库对象
var db *gorm.DB

// 数据库初始化函数
func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.DataBaseSettings.User,
		settings.DataBaseSettings.Password,
		settings.DataBaseSettings.Host,
		settings.DataBaseSettings.Name)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{PrepareStmt: true})
	if err != nil {
		log.Fatal(2, "连接数据库失败: %v", err)
	}
	logging.Info("数据库连接成功")
}
