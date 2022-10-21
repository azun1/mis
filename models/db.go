package models

import (
	"MIS/pkg/logging"
	"MIS/pkg/settings"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"time"
)

var db *gorm.DB
var err error

// 连接配置数据库
func initDB() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		settings.DataBaseSettings.User,
		settings.DataBaseSettings.Password,
		settings.DataBaseSettings.Host,
		settings.DataBaseSettings.Name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
		},
		PrepareStmt: true})
	if err != nil {
		log.Fatal(2, "连接数据库失败: %v", err)
	}
	logging.Info("数据库连接成功")

	sqlDB, err := db.DB()
	// SetMaxIdleConns 设置空闲连接池中连接的最大数量
	sqlDB.SetMaxIdleConns(10)

	// SetMaxOpenConns 设置打开数据库连接的最大数量。
	sqlDB.SetMaxOpenConns(100)

	// SetConnMaxLifetime 设置了连接可复用的最大时间。
	sqlDB.SetConnMaxLifetime(10 * time.Second)
}
