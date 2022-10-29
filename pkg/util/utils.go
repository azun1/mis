package util

import (
	"MIS/pkg/logging"
	"github.com/google/uuid"
	"strconv"
	"strings"
)

// GenerateUuid :生成UUID
func GenerateUuid() string {
	Uuid := uuid.New()
	key := Uuid.String()
	logging.Info("GenerateUuid key: %s\n", key)
	return key
}

type YMD struct {
	Year  int
	Month int
	Day   int
}

// DateStringToYMD 日期字符串转年月日
func DateStringToYMD(date string) (ymd YMD, err error) {
	a := strings.Split(date, "-")
	ymd.Year, err = strconv.Atoi(a[0])
	if err != nil {
		return
	}
	ymd.Month, err = strconv.Atoi(a[1])
	if err != nil {
		return
	}
	ymd.Day, err = strconv.Atoi(a[2])
	if err != nil {
		return
	}

	return
}
