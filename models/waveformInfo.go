package models

import (
	"MIS/pkg/logging"
	"encoding/csv"
	"errors"
	"gorm.io/gorm"
	"io"
	"os"
	"sort"
	"strings"
	"time"
)

type WaveformInfo struct {
	gorm.Model
	// 联合索引遵循左前缀原则, priority 数值小的在左边(默认值为10), 相同则按在 struct 中的次序排
	UserUuid  string    `json:"user_uuid" gorm:"default:'';not null;uniqueIndex:u&f_index,priority:1;uniqueIndex:u&w&f_index,priority:1"`
	StartTime time.Time `json:"start_time" gorm:"type:timestamp;default:'2000-01-01 00:00:01';not null"`
	EndTime   time.Time `json:"end_time" gorm:"type:timestamp;default:'2000-01-01 00:00:01';not null"`
	WaveType  string    `json:"wave_type" gorm:"type:varchar(50);default:'';not null;uniqueIndex:u&w&f_index,priority:2;uniqueIndex:w&f_index,priority:1"`
	FilePath  string    `json:"file_path" gorm:"type:varchar(500);default:'';not null;uniqueIndex:u&f_index,priority:2;uniqueIndex:u&w&f_index,priority:3;uniqueIndex:w&f_index,priority:2"`
}

type Description struct {
	WaveformType string   `json:"waveform_type"`
	TimeStart    string   `json:"time_start"`
	TimeEnd      string   `json:"time_end"`
	ValuesMin    string   `json:"values_min"`
	ValuesMax    string   `json:"values_max"`
	X            []string `json:"x"`
	Y            []string `json:"y"`
}

// CreateNewRecord 对外接口, 传感器上传数据的时候触发
// TODO: 硬件那边最好可以加个蓝牙
func CreateNewRecord() error {

	return nil
}

// ReadByLine 按行导出n行数据(界面展示, 指出危险波形)
// @params: fp是文件路径;  n为可变参数, 若调用时省略n则默认输出最新的一行数据
func ReadByLine(fp string, data *[][]string, n ...int) error {
	file, err := os.Open(fp)
	if err != nil {
		return err
	}

	defer file.Close()

	reader := csv.NewReader(file)
	// 设置每行数据里面的字段数量
	// 如果FieldsPerRecord的值为正数, 并且CSV文件里面读取出的字段数量少于这个值时, Go就会抛出一个错误
	// 如果FieldsPerRecord的值为0, 那么读取器就会将读取到的第一条记录的字段数量用作FieldsPerRecord的值
	// 如果FieldsPerRecord的值为负数, 允许读取到的字段数量为任意值
	reader.FieldsPerRecord = -1
	reader.ReuseRecord = true // 允许复用切片
	var record []string
	var line int
	var size = len(n)

	if size > 0 {
		// 先排个序
		// Ref: https://cloud.tencent.com/developer/article/1071728
		if !sort.IntsAreSorted(n) {
			sort.Ints(n)
		}

		// 导出指定行(从0开始计数, 第0行是数据元信息行)的数据
		for i := 0; i < size; {
			record, err = reader.Read()
			if err == io.EOF {
				break
			}
			if line == n[i] {
				// 追加的内容形式: [0.9915713654170097 1.0 0.0]
				csvData := make([]string, len(record))
				copy(csvData, record)
				*data = append(*data, csvData)
				i++
			}
			line++
		}
	} else {
		// 读取最后一行
		// TODO: 这里肯定有比顺序查找更好的算法, 但是我看了一圈都很复杂, 以后再说
		var csvData []string

		for {
			record, err = reader.Read() // 按行读取数据,可控制读取部分
			if err == io.EOF {
				// 读取到文件尾时, 返回的record为nil
				break
			}
			csvData = record // 引用赋值
		}
		// 忽略行首序号
		*data = append(*data, csvData[1:])
	}
	return nil
}

// GetLatestRate 通过发送包含数据描述的波形JSON, <对接>前端波形展示界面
func (u *User) GetLatestRate(desc *Description) (err error) {
	wfi := WaveformInfo{}
	// 找到最近一个记录
	res := db.Where("user_uuid = ? AND wave_type = ?", u.Uuid, desc.WaveformType).
		Order("start_time desc").First(&wfi)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		logging.Info("There is no waveform record for now")
		return gorm.ErrRecordNotFound
	}
	var et = wfi.EndTime
	// csv文件一行10s
	st := et.Add(-time.Second * 10) // .Sub方法只能两个Time相减, .Add 才能加减时长
	// time.Time 转 字符串固定写法
	desc.TimeStart = st.Format("2006-01-02 15:04:05")
	desc.TimeEnd = et.Format("2006-01-02 15:04:05")
	var data [][]string
	err = ReadByLine(wfi.FilePath, &data)
	if err != nil {
		return err
	}
	err = RowsToJSON(desc, &data)
	if err != nil {
		return err
	}
	return nil
}

// RowsToCSV 以行为单位, 将数据转成前端展示需要的CSV
// 发送给前端CSV文件格式: X轴时序, Y轴数据值
func RowsToCSV(des Description, data *[][]string) error {

	return nil
}

// RowsToJSON 以行为单位, 将数据序列化成前端展示需要的JSON
// 发送给前端JSON文件格式: tx字节切片 - X轴时序, ty字节切片 - Y轴数据值
func RowsToJSON(desc *Description, data *[][]string) (err error) {
	st, _ := time.Parse("2006-01-02 15:04:05", desc.TimeStart)
	et, _ := time.Parse("2006-01-02 15:04:05", desc.TimeEnd)

	// X轴处理
	dur := et.Sub(st)
	// Duration类型只能和int64, float64, string转换
	var xLng int64 = int64(len((*data)[0]))
	// 最小时隙(用int6[]string
	gap := int64(dur) / xLng
	x := make([]time.Time, 0, xLng)
	desc.X = make([]string, 0, xLng)
	for i := 0; i < cap(x); i++ {
		x = append(x, st)
		x[i] = x[i].Add(time.Duration(int64(i) * gap))
		// .000 精确到毫秒
		desc.X = append(desc.X, x[i].Format("2006-01-02 15:04:05.000"))
	}

	// Y轴处理
	var yLng int
	for i := range *data {
		yLng += len((*data)[i])
	}
	desc.Y = make([]string, 0, yLng)
	for i := range *data {
		desc.Y = append(desc.Y, (*data)[i]...)
	}

	// 找出Y轴最值
	maxV := desc.Y[0]
	minV := desc.Y[1]
	for i := range desc.Y {
		if strings.Compare(desc.Y[i], maxV) > 0 {
			maxV = desc.Y[i]
		}
		if strings.Compare(desc.Y[i], minV) < 0 {
			minV = desc.Y[i]
		}
	}

	desc.ValuesMax = maxV
	desc.ValuesMin = minV

	return err
}

// DownloadCSVByType 按类型导出CSV文件
// 按类型导出该用户的所有原始CSV数据给AI分析
// @params: fps是文件路径切片
func (u *User) DownloadCSVByType(WaveformType string, fps *[]string) error {
	var wfi []WaveformInfo
	res := db.Where("user_uuid = ? AND wave_type = ?", u.Uuid, WaveformType).
		Order("start_time desc").Find(&wfi)
	if len(wfi) == 0 {
		logging.Info("There is no such type waveform record for now")
		return gorm.ErrRecordNotFound
	}
	for i := range wfi {
		*fps = append(*fps, wfi[i].FilePath)
	}
	return res.Error
}
