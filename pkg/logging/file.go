package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/" // 文件路径
	LogSaveName = "log"           // 文件名
	LogFileExt  = "log"           // 文件后缀
	TimeFormat  = "20060102"
)

func getLogFilePath() string {
	return fmt.Sprintf("%s", LogSavePath)
}

func getLogFileFullPath() string {
	prefixPath := getLogFilePath()
	suffixPath := fmt.Sprintf("%s%s.%s", LogSaveName, time.Now().Format(TimeFormat), LogFileExt)

	return fmt.Sprintf("%s%s", prefixPath, suffixPath)
}

func openLogFile(filePath string) *os.File {
	_, err := os.Stat(filePath) // 返回文件信息结构描述文件
	switch {
	case os.IsNotExist(err): // 判断文件是否存在
		mkDir()
	case os.IsPermission(err): // 判断权限
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}

// 新建目录，未新建文件
func mkDir() {
	dir, _ := os.Getwd() // 返回与当前目录对应的根路径名（项目路径）
	// 创建对应的目录以及所需的子目录，若成功则返回nil，否则返回error
	err := os.MkdirAll(dir+"/"+getLogFilePath(), os.ModePerm) // const定义ModePerm FileMode = 0777
	if err != nil {
		panic(err)
	}
}
