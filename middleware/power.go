package middleware

import (
	"MIS/models"
	"MIS/pkg/settings"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Can api权限设置中间件
func Can(power ...settings.UserType) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := c.MustGet("user").(*models.User)
		match := false

		for k := range power {
			if power[k].Int() == user.Power {
				match = true
				break
			}
		}
		if !match {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "forbidden",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
