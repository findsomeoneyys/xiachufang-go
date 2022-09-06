package router

import (
	docs "github.com/findsomeoneyys/xiachufang-api/docs"
	"github.com/findsomeoneyys/xiachufang-api/router/api"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter initialize routing information
func InitRouter(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	apiGroup := r.Group("/api")

	//健康检查
	apiGroup.HEAD("/healthcheck", api.HealthCheck)

	//搜索
	apiGroup.GET("/search/:keyword", api.Search)

	//分类
	apiGroup.GET("/category", api.GetAllCategory)
	apiGroup.GET("/category/:no", api.SearchCategory)

	// 菜谱
	apiGroup.GET("/recipe/:no", api.GetRecipe)
}
