package models

import (
	"MIS/pkg/logging"
	"MIS/pkg/settings"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
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
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   settings.DataBaseSettings.TablePrefix, // 设置表前缀
			SingularTable: true,                                  // 禁用默认表名的复数形式，如果置为 true，则 `Common` 的默认表名是 `user`
		},
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatalf("连接数据库失败: %v\n", err)
	}
	logging.Info("数据库连接成功")

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&MessageRecord{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&SystemNotice{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&UserRelationship{})
	if err != nil {
		log.Fatalln(err)
	}
	err = db.AutoMigrate(&WaveformInfo{})
	if err != nil {
		log.Fatalln(err)
	}
}
