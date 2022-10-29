package v1

import (
	"MIS/api"
	"MIS/models"
	"MIS/pkg/settings"
	"MIS/pkg/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// Register 用户注册
func Register(context *gin.Context) {
	var user models.User
	if !util.BindAndValid(context, &user) {
		return
	}

	// 判断用户名是否已存在
	if models.IsUserNameExist(user.Name) {
		code := -1
		context.JSON(code, gin.H{
			"code":    code,
			"message": "用户名已存在",
		})
		return
	}
	// 生成uuid
	user.Uuid = util.GenerateUuid()
	// 设置用户类型
	user.Power = settings.Common.Int()
	// 在数据库中添加用户
	err := user.Add()
	if err != nil {
		api.ErrHandle(context, err)
		return
	}
	// 颁发token
	token, err := util.GenerateToken(user.Uuid)
	if err != nil {
		api.ErrHandle(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"mesg":  "注册成功",
		"token": token,
	})
}

// Login 用户登录
func Login(context *gin.Context) {
	var userLoginForm struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
	if !util.BindAndValid(context, &userLoginForm) {
		return
	}
	// 查找用户
	user, err := models.FindUserByName(userLoginForm.Username)
	if err != nil {
		code := -1
		context.JSON(code, gin.H{
			"code":    code,
			"message": "该用户不存在",
		})
		return
	}
	// 比对密码
	if user.Password != userLoginForm.Password {
		code := -1
		context.JSON(code, gin.H{
			"code":    code,
			"message": "密码错误",
		})
		return
	}
	// 颁发token
	token, err := util.GenerateToken(user.Uuid)
	if err != nil {
		api.ErrHandle(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code":  http.StatusOK,
		"mesg":  "登录成功",
		"token": token,
	})
}

// Logout 用户登出
func Logout(context *gin.Context) {
	user := api.CurrentUser(context)
	// 更新用户上次登出时间，使当前token失效
	err := user.UpdateUserTime()
	if err != nil {
		api.ErrHandle(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "登出成功",
	})
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(context *gin.Context) {
	var userForm struct {
		Username string `json:"username"`
		RealName string `json:"realName"`
		Email    string `json:"email" validate:"omitempty,email"`
		Gender   string `json:"gender"`
		Birth    string `json:"birth" validate:"omitempty,datetime=2006-01-02"`
	}

	if !util.BindAndValid(context, &userForm) {
		return
	}

	birthPointer := new(time.Time)
	if userForm.Birth != "" {
		ymd, err := util.DateStringToYMD(userForm.Birth)
		if err != nil {
			api.ErrHandle(context, err)
			return
		}
		birth := time.Date(ymd.Year, time.Month(ymd.Month), ymd.Day, 0, 0, 0, 0, time.UTC)
		*birthPointer = birth
	} else {
		birthPointer = nil
	}

	n := models.User{
		Email:    userForm.Email,
		Name:     userForm.Username,
		RealName: userForm.RealName,
		Gender:   userForm.Gender,
		Birth:    birthPointer,
	}

	user := api.CurrentUser(context)
	err := user.Update(n)
	if err != nil {
		api.ErrHandle(context, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "用户信息更新成功",
	})
}

// Delete 注销用户
func Delete(context *gin.Context) {
	user := api.CurrentUser(context)
	err := user.Delete()
	if err != nil {
		api.ErrHandle(context, err)
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"mesg": "注销成功",
	})
}
