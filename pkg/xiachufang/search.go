package xiachufang

import (
	"log"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) Search(keyword string, page int) (*RecipeListResult, error) {
	u, _ := GetApiURL("SEARCH")

	queryString := u.Query()
	queryString.Set("keyword", keyword)
	queryString.Set("page", strconv.Itoa(page))
	u.RawQuery = queryString.Encode()

	doc, err := c.Visit(u.String())
	if err != nil {
		log.Fatal("client.Search.Visit error : ", err)
		return nil, err
	}

	res := ParseSearchPage(doc)
	return res, nil

}

func ParseSearchPage(doc *goquery.Document) *RecipeListResult {
	recipes := make([]*Recipe, 0)

	// 解析搜索页菜谱列表信息
	doc.Find(".normal-recipe-list .recipe").Each(func(i int, s *goquery.Selection) {
		// 封面
		imgCover := s.Find("a").First()

		link, exit := imgCover.Attr("href")
		if exit {
			link, _ = UrlRelativeToAbsolute(link)
		}

		imgSrc, _ := imgCover.Find(".cover img").Attr("data-src")

		// 详情
		info := s.Find(".info")

		name := info.Find(".name").Text()
		name = strings.TrimSpace(name)

		// 用料  搜索页结果 只有名字，没有用量
		materials := make([]*RecipeMaterial, 0)
		info.Find(".ing").Children().Each(func(b int, ing *goquery.Selection) {
			materialName := ing.Text()
			materialLink, exit := ing.Attr("href")
			if exit {
				materialLink, _ = UrlRelativeToAbsolute(materialLink)
			}
			materials = append(materials, &RecipeMaterial{Name: materialName, Link: materialLink})
		})

		// 评分
		stats := s.Find(".stats")
		var score float64
		var cooked int

		stats.Find("span").Each(func(i int, s *goquery.Selection) {
			if s.HasClass("green-font") {
				score, _ = strconv.ParseFloat(s.Text(), 64)
			} else {
				cooked, _ = strconv.Atoi(s.Text())
			}
		})

		// 作者
		authorNode := s.Find(".author a")
		authorName := authorNode.Text()
		authorLink, ok := authorNode.Attr("href")
		if ok {
			authorLink, _ = UrlRelativeToAbsolute(authorLink)
		}
		authors := make([]*Author, 0)
		authors = append(authors, &Author{
			Name: authorName,
			Link: authorLink,
		})

		recipe := &Recipe{
			Name:      name,
			Link:      link,
			Cover:     imgSrc,
			Authors:   authors,
			Score:     score,
			Cooked:    cooked,
			Materials: materials,
		}
		recipes = append(recipes, recipe)
	})

	//解析当前页， 上/下页
	var prev string
	var current int
	var next string

	prevLink, ok := doc.Find(".pager a.prev").First().Attr("href")
	if ok {
		prev, _ = UrlRelativeToAbsolute(prevLink)
	}

	now, err := strconv.Atoi(doc.Find(".pager .now").First().Text())
	if err == nil {
		current = now
	}

	nextLink, ok := doc.Find(".pager a.next").First().Attr("href")
	if ok {
		next, _ = UrlRelativeToAbsolute(nextLink)
	}

	res := &RecipeListResult{
		Recipes:  recipes,
		PrevPage: prev,
		Current:  current,
		NextPage: next,
	}
	return res
}
