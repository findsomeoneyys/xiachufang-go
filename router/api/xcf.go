package api

import (
	"net/http"
	"strconv"

	"github.com/findsomeoneyys/xiachufang-api/pkg/app"
	"github.com/findsomeoneyys/xiachufang-api/pkg/code"
	"github.com/findsomeoneyys/xiachufang-api/pkg/xiachufang"
	"github.com/gin-gonic/gin"
)

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

// @Description 下厨房搜索 可搜索分类/菜谱..
// @Tags 下厨房
// @Accept  json
// @Produce  json
// @Param   keyword     path    string     true        "搜索关键字"
// @Param   page      query    int     false        "获取第几页, 默认第一页"
// @Success 200 {object} app.Response{data=xiachufang.RecipeListResult} "{"code": 200, "data": [...]}"
// @Router /api/search/{keyword} [get]
func Search(c *gin.Context) {
	appG := app.Gin{C: c}
	client := xiachufang.NewClient()

	keyword := c.Param("keyword")
	page, e := strconv.Atoi(c.DefaultQuery("page", "1"))

	if e != nil {
		appG.Response(http.StatusBadRequest, code.INVALID_PARAMS, e.Error(), nil)
		return
	}

	res, err := client.Search(keyword, page)

	if err != nil {
		appG.Response(http.StatusInternalServerError, code.ERROR, err.Error(), nil)
		return
	}

	appG.Response(http.StatusOK, code.SUCCESS, "", res)

}

// @Description 获取下厨房全部分类
// @Tags 下厨房
// @Accept  json
// @Produce  json
// @Success 200 {object} app.Response{data=map[string][]xiachufang.Category} "{"code": 200, "data": [...]}"
// @Router /api/category/ [get]
func GetAllCategory(c *gin.Context) {
	appG := app.Gin{C: c}
	client := xiachufang.NewClient()

	res, err := client.GetAllCategory()

	if err != nil {
		appG.Response(http.StatusInternalServerError, code.ERROR, err.Error(), nil)
		return
	}

	appG.Response(http.StatusOK, code.SUCCESS, "", res)

}

// @Description 获取下厨房某个分类菜谱
// @Tags 下厨房
// @Accept  json
// @Produce  json
// @Param  no    path  string  true  "分类"
// @Param  sort    query  string  false  "排序方式 默认:最近流行 pop:最受欢迎 time:评分"
// @Param  page      query    int     false        "获取第几页, 默认第一页"
// @Success 200 {object} app.Response{data=map[string][]xiachufang.Category} "{"code": 200, "data": [...]}"
// @Router /api/category/{no} [get]
func SearchCategory(c *gin.Context) {
	appG := app.Gin{C: c}
	client := xiachufang.NewClient()

	no := c.Param("no")

	// sort排序 默认最新流行 pop最受欢迎 time 评分
	var searchType xiachufang.SearchCategoryType
	s := c.DefaultQuery("sort", "")
	switch s {
	case "pop":
		searchType = xiachufang.SearchCategoryTypePopular
	case "time":
		searchType = xiachufang.SearchCategoryTypeTime
	default:
		searchType = xiachufang.SearchCategoryTypeRecent
	}

	page, e := strconv.Atoi(c.DefaultQuery("page", "1"))

	if e != nil {
		appG.Response(http.StatusBadRequest, code.INVALID_PARAMS, e.Error(), nil)
		return
	}

	res, err := client.SearchCategory(no, searchType, page)

	if err != nil {
		appG.Response(http.StatusInternalServerError, code.ERROR, err.Error(), nil)
		return
	}

	appG.Response(http.StatusOK, code.SUCCESS, "", res)

}

// @Description 下厨房 获取菜谱
// @Tags 下厨房
// @Accept  json
// @Produce  json
// @Param   no     path    string     true        "菜谱编号"
// @Success 200 {object} app.Response{data=xiachufang.Recipe} "{"code": 200, "data": [...]}"
// @Router /api/recipe/{no} [get]
func GetRecipe(c *gin.Context) {
	appG := app.Gin{C: c}
	client := xiachufang.NewClient()

	no := c.Param("no")
	res, err := client.GetRecipe(no)

	if err != nil {
		appG.Response(http.StatusInternalServerError, code.ERROR, err.Error(), nil)
		return
	}

	appG.Response(http.StatusOK, code.SUCCESS, "", res)

}
