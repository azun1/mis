package settings

import (
	"github.com/go-ini/ini"
	"log"
	"os"
)

// app.ini配置文件读取初始化模块

var Conf *ini.File

type UserType int // 用户类型

const (
	None UserType = iota
	Common
	Administrator
)

func (u UserType) Int() int {
	return int(u)
}

type App struct {
	PageSize  int
	JwtSecret string
}

type Server struct {
	HttpPort string
	RunMode  string
}

type DataBase struct {
	User        string
	Password    string
	Host        string
	Name        string
	TablePrefix string
}

type OssSetting struct {
	Endpoint         string
	ArchiveDirectory string
}

type AliyunSetting struct {
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string
	BucketUrl       string
}

var AppSettings = &App{}
var ServerSettings = &Server{}
var DataBaseSettings = &DataBase{}
var OssSettings = &OssSetting{}
var AliyunSettings = &AliyunSetting{}

func init() {
	var err error
	Conf, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.go, 读取项目配置文件 'app.ini'失败: %v", err)
	}

	mapTo("app", AppSettings)
	mapTo("server", ServerSettings)
	mapTo("database", DataBaseSettings)
	mapTo("ossSetting", OssSettings)
	mapTo("aliyunSetting", AliyunSettings)

	// 创建存储文件夹
	_, err = os.Stat(OssSettings.ArchiveDirectory)
	if os.IsNotExist(err) {
		// ModePerm: 最高权限，可读可写可执行
		err = os.MkdirAll(OssSettings.ArchiveDirectory, os.ModePerm)
	}
}

func mapTo(section string, v interface{}) {
	err := Conf.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("Cfg.MapTo %s err: %v", section, err)
	}
}
