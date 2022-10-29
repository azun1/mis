package middleware

import (
	"MIS/api"
	"MIS/models"
	"MIS/pkg/e"
	"MIS/pkg/util"
	"errors"
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
		token := c.Request.Header.Get("token")
		if token == "" {
			code = e.INVALID_PARAMS
		} else {
			claims, err := util.ParseToken(token)
			if err != nil {
				code = e.ERROR_AUTH_CHECK_TOKEN_FAIL
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = e.ERROR_AUTH_CHECK_TOKEN_TIMEOUT
			} else {
				// 获取用户信息
				user, err := models.FindUserByUuid(claims.UserUuid)
				if err != nil {
					api.ErrHandle(c, err)
					c.Abort()
					return
				}
				// 检查该token是否已登出
				if user.LogoutTime != nil && user.LogoutTime.Before(time.Unix(claims.IssuedAt, 0)) {
					api.ErrHandle(c, errors.New("当前token已登出，请重新登录"))
					c.Abort()
					return
				}
				// 更新用户上次访问时间
				err = user.UpdateLastActiveTime()
				if err != nil {
					api.ErrHandle(c, err)
					c.Abort()
					return
				}
				c.Set("user", &user) // 后续接口处理可直接通过c.Get(key)获取用户信息
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
