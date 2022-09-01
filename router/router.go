package router

import (
	"github.com/findsomeoneyys/xiachufang-api/router/api"
	"github.com/gin-gonic/gin"
)

// InitRouter initialize routing information
func InitRouter(r *gin.Engine) {

	group := r.Group("/api")
	group.GET("/search/:keyword", api.Search)

	//分类
	group.GET("/category", api.GetAllCategory)
	group.GET("/category/:no/*searchType", api.SearchCategory)

	// 菜谱
	group.GET("/recipe/:no", api.GetRecipe)
}
