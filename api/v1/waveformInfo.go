package v1

import (
	"MIS/api"
	"MIS/models"
	"archive/zip"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

// 临时存放压缩文件的目录
var archiveDirectory = "/home/healthy_bigdata/data/temp/"

// AllMineHeartRate 获得自己所有的心率数据
func AllMineHeartRate(c *gin.Context) {
	var user = api.CurrentUser(c)
	var filePaths []string
	err := user.DownloadCSVByType("HeartRate", &filePaths)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	// 设置压缩文件存放路径(以用户名称命名)
	var archivePath = archiveDirectory + user.Name + "HeartRate.zip"

	err = Compress(&filePaths, archivePath)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	// 发送压缩文件
	c.FileAttachment(archivePath, "HeartRate.zip")

}

// AllMineBreathRate 获得自己的所有呼吸率数据
func AllMineBreathRate(c *gin.Context) {
	var user = api.CurrentUser(c)
	var filePaths []string
	err := user.DownloadCSVByType("BreathRate", &filePaths)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	// 设置压缩文件存放路径(以用户名称命名)
	var archivePath = archiveDirectory + user.Name + "BreathRate.zip"

	err = Compress(&filePaths, archivePath)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	// 发送压缩文件
	c.FileAttachment(archivePath, "BreathRate.zip")
}

// LatestHeartRate 获得最近的心率数据 (10s)
func LatestHeartRate(c *gin.Context) {
	var user = api.CurrentUser(c)
	var desc = models.Description{
		WaveformType: "HeartRate",
	}
	err := user.GetLatestRate(&desc)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, desc)
}

// LatestBreathRate 获得最近的呼吸率数据 (10s)
func LatestBreathRate(c *gin.Context) {
	var user = api.CurrentUser(c)
	var desc = models.Description{
		WaveformType: "BreathRate",
	}
	err := user.GetLatestRate(&desc)
	if err != nil {
		api.ErrHandle(c, err)
		return
	}
	c.JSON(http.StatusOK, desc)
}

// Compress 压缩文件
// @src: source, 文件路径数组 @dst: destination, 压缩archive存放位置
func Compress(src *[]string, dst string) (err error) {
	// 1. 创建zip archive文件(一个普通的文件)
	archive, err := os.Create(dst)
	if err != nil {
		return err
	}
	// 在最后关闭archive文件
	defer func() {
		if err = archive.Close(); err != nil {
			log.Fatalln(err)
			// 这里不要return err, 会提示『要返回的实参太多』
			// https://www.cnblogs.com/phpper/archive/2022/06/18/16389393.html#%E7%AC%94%E8%AF%95%E9%A2%98%E4%B8%80
		}
	}()
	// 2. 创建zipWriter写入流
	zipWriter := zip.NewWriter(archive)

	// 3. 添加文件: 首先, 获取文件内容; 然后使用zw.Create()来指定我们想要保存该文件的位置
	for i := range *src {
		f, _ := os.Open((*src)[i])
		err = compress(f, strconv.Itoa(i+1), zipWriter)
		if err != nil {
			log.Fatalln(err)
		}
	}
	// 4. 关闭zipWriter写入流, 所有文件都会被逐一压缩并被保存在archive中
	if err = zipWriter.Close(); err != nil {
		log.Fatalln(err)
		return err
	}

	// 上面defer调用的匿名函数可以给err赋值
	return err
}

func compress(file *os.File, prefix string, zw *zip.Writer) (err error) {
	// 因为循环里面不能defer, 我才写了这个函数
	defer file.Close()

	// 使用Create()添加条目
	w, err := zw.Create(prefix + ".csv")
	if err != nil {
		return err
	}
	// 复制文件内容到对应条目
	if _, err := io.Copy(w, file); err != nil {
		return err
	}
	return nil
}
