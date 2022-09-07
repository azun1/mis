package routers

import (
	"MIS/middleware/jwt"
	"MIS/pkg/settings"
	"MIS/routers/api"
	v1 "MIS/routers/api/v1"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(settings.RunMode)

	r.GET("/auth", api.GetAuth)

	apiv1 := r.Group("/api/v1")
	apiv1.Use(jwt.JWT()) // token校验
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//获取商品列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定商品
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建商品
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定商品
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定商品
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
