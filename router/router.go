package router

import (
	"github.com/findsomeoneyys/xiachufang-api/router/api"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter(r *gin.Engine) {

	apiGroup := r.Group("/api")

	//健康检查
	apiGroup.HEAD("/healthcheck", api.HealthCheck)

	//搜索
	apiGroup.GET("/search/:keyword", api.Search)

	//分类
	apiGroup.GET("/category", api.GetAllCategory)
	apiGroup.GET("/category/:no/*searchType", api.SearchCategory)

	// 菜谱
	apiGroup.GET("/recipe/:no", api.GetRecipe)
}
