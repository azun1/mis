package util

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"log"
	"net/http"
	"strings"
)

// BindAndValid 表单数据验证函数
func BindAndValid(context *gin.Context, data interface{}) bool {
	// 实例化验证对象
	var validate = validator.New()
	err := context.ShouldBind(data)
	if err != nil {
		log.Println(err.Error())
		context.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"err":  err.Error(),
		})
		return false
	}
	err = validate.Struct(data)
	if err != nil {
		log.Println(err)
		trans := validateTransInit(validate)
		verrs := err.(validator.ValidationErrors)
		errs := make(map[string]string)
		for key, value := range verrs.Translate(trans) {
			errs[key[strings.Index(key, ".")+1:]] = value
		}
		context.JSON(http.StatusNotAcceptable, gin.H{
			"errors": errs,
			"mesg":   "请求参数错误",
			"code":   http.StatusNotAcceptable,
		})
		return false
	}
	return true
}

// validateTransInit 表单数据验证翻译器
func validateTransInit(validate *validator.Validate) ut.Translator {
	// 万能翻译器，保存所有的语言环境和翻译数据
	uni := ut.New(zh.New())
	// 翻译器
	trans, _ := uni.GetTranslator("zh")
	//验证器注册翻译器
	err := zhTranslations.RegisterDefaultTranslations(validate, trans)
	if err != nil {
		fmt.Println(err)
	}
	return trans
}
