package v1

import "github.com/gin-gonic/gin"

// RequestConnect 发送(1条)关联账号请求
func RequestConnect(c *gin.Context) {
	// 根据输入的email / username 查找对应用户
	//var userSelectForm struct {
	//	Username string `json:"username" validate:"required"`
	//}
	// 如何确定用户输入的是name还是email, 用反射来做判断吗
}

// AcceptConnect 同意(1条)关联账号请求
func AcceptConnect(c *gin.Context) {

}

// DeleteConnection 拒绝/删除(1条)关联关系
func DeleteConnection(c *gin.Context) {

}

// GetRelatedAccList 获取(同意/未同意的)关联账号列表
func GetRelatedAccList(c *gin.Context) {

}

// GetRelatedAccount 获取某个已关联账号的信息(关系类型, 备注)
func GetRelatedAccount(c *gin.Context) {

}

// SetRelatedAccount 设置某个已关联账号的信息(关系类型, 备注)
func SetRelatedAccount(c *gin.Context) {

}
