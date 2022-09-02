package api

import (
	"net/http"
	"strconv"

	"github.com/findsomeoneyys/xiachufang-api/pkg/app"
	"github.com/findsomeoneyys/xiachufang-api/pkg/code"
	"github.com/findsomeoneyys/xiachufang-api/pkg/xiachufang"
	"github.com/gin-gonic/gin"
)

// res, _ := client.Search("西瓜", 2)
// res, _ := client.GetRecipe("101829462")

func HealthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
	})
}

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

func SearchCategory(c *gin.Context) {
	appG := app.Gin{C: c}
	client := xiachufang.NewClient()

	no := c.Param("no")

	// searchType排序 默认最新流行 pop最受欢迎 time 评分
	var searchType xiachufang.SearchCategoryType
	s := c.Param("searchType")
	switch s {
	case "/pop":
		searchType = xiachufang.SearchCategoryTypePopular
	case "/time":
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
