package jwt

import (
	"MIS/pkg/e"
	"MIS/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// JWT jwt验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}

		code = e.SUCCESS
		token := c.Query("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			} else {
				c.Set("username", claims.Username) // 后续接口处理可直接通过c.Get(key)获取用户信息
			}
		}

		if code != e.SUCCESS {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": code,
				"msg":  e.GetMsg(code),
				"data": data,
			})

			c.Abort() // 终止运行后续函数
			return
		}

		c.Next()
	}
}
