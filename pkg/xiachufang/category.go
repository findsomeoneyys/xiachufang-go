package xiachufang

import (
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 搜索分类时可选排序方式
type SearchCategoryType string

const (
	SearchCategoryTypeRecent  SearchCategoryType = ""     //最近流行 默认值
	SearchCategoryTypePopular SearchCategoryType = "pop"  //最受欢迎
	SearchCategoryTypeTime    SearchCategoryType = "time" //评分
)

func (c *Client) GetAllCategory() (map[string][]*Category, error) {
	u, _ := GetApiURL("CATEGORY")

	doc, err := c.Visit(u.String())
	if err != nil {
		log.Fatal("client.Search.Visit error : ", err)
		return nil, err
	}

	res := make(map[string][]*Category, 0)

	base := u.String()
	doc.Find(".category-container .cates-list").Each(func(i int, s *goquery.Selection) {
		// 最外层大分类名 eg 热门专题/烘焙甜品饮料/肉类..
		topCateName := strings.TrimSpace(s.Find(".cates-list-info h3").Text())

		s.Find(".cates-list-all h4").Each(func(i int, h *goquery.Selection) {
			// 大分类里每一个小分类 比如 菜式，特色食品，..
			subCateName := strings.TrimSpace(h.Text())
			subCateNo, exits := h.Attr("id")
			if exits {
				subCateNo = strings.Trim(subCateNo, "cat")
			}
			cate := Category{
				No:       subCateNo,
				Name:     subCateName,
				Link:     fmt.Sprintf("%s/%s/", base, subCateNo),
				Children: make([]*Category, 0),
			}

			// 取小分类里面全部
			h.Next().Find("li").Each(func(i int, li *goquery.Selection) {
				sNo, exits := li.Attr("id")
				if exits {
					sNo = strings.Trim(sNo, "cat")
				}

				aNode := li.Find("a").First()
				name := strings.TrimSpace(aNode.Text())
				link, exits := aNode.Attr("href")
				if exits {
					link, _ = UrlRelativeToAbsolute(link)
				}

				c := &Category{
					No:   sNo,
					Name: name,
					Link: link,
				}
				cate.Children = append(cate.Children, c)
			})

			res[topCateName] = append(res[topCateName], &cate)
		})

	})

	return res, nil

}

// 搜索分类需要提供排序, 默认页为SearchCategoryTypeRecent
func (c *Client) SearchCategory(no string, searchType SearchCategoryType, page int) (*RecipeListResult, error) {
	u, _ := GetApiURL("CATEGORY")

	u.Path = path.Join(u.Path, fmt.Sprintf("%s/", no))

	switch searchType {
	case SearchCategoryTypeRecent:
	case SearchCategoryTypePopular:
		u.Path = path.Join(u.Path, "/pop/")
	case SearchCategoryTypeTime:
		u.Path = path.Join(u.Path, "/time/")
	}

	queryString := u.Query()
	queryString.Set("page", strconv.Itoa(page))
	u.RawQuery = queryString.Encode()

	doc, err := c.Visit(u.String())
	if err != nil {
		log.Fatal("client.Search.Visit error : ", err)
		return nil, err
	}

	res := ParseCategory(doc)
	return res, nil

}

func ParseCategory(doc *goquery.Document) *RecipeListResult {
	recipes := make([]*Recipe, 0)
	fmt.Println(doc.Find(".main-panel").Size())

	// 解析搜索页菜谱列表信息
	doc.Find(".main-panel .normal-recipe-list div.recipe").Each(func(i int, s *goquery.Selection) {
		// 封面
		imgCover := s.Find("a").First()

		link, exit := imgCover.Attr("href")
		if exit {
			link, _ = UrlRelativeToAbsolute(link)
		}

		imgSrc, exits := imgCover.Find(".cover img").Attr("data-src")
		if exits {
			imgSrc = strings.Split(imgSrc, "?")[0]
		}

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
			Name:           name,
			Link:           link,
			Cover:          imgSrc,
			Authors:        authors,
			Score:          score,
			CookedLast7Day: cooked,
			Materials:      materials,
			Steps:          make([]*RecipeStep, 0),
			Tip:            "",
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
