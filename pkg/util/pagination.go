package util

import (
	"MIS/pkg/settings"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage 返回请求页码对应记录起始位置
func GetPage(c *gin.Context) int {
	result := 0
	page, _ := com.StrTo(c.Query("page")).Int()
	if page > 0 {
		result = (page - 1) * settings.AppSettings.PageSize
	}

	return result
}

// TotalPage 返回总记录数可分为多少页
func TotalPage(total int64) int64 {
	n := total / int64(settings.AppSettings.PageSize)
	if total%int64(settings.AppSettings.PageSize) > 0 {
		n++
	}
	return n
}

// GetPagination 返回分页相关返回信息
func GetPagination(c *gin.Context, total int64, currentPage int) gin.H {
	return gin.H{
		"total":        total,
		"per_page":     settings.AppSettings.PageSize,
		"current_page": currentPage,
		"total_pages":  TotalPage(total),
	}
}
