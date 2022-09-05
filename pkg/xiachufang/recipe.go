package xiachufang

import (
	"fmt"
	"log"
	"path"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (c *Client) GetRecipe(no string) (*Recipe, error) {
	u, _ := GetApiURL("RECIPE")

	u.Path = path.Join(u.Path, fmt.Sprintf("%s/", no))

	doc, err := c.Visit(u.String())
	if err != nil {
		log.Fatal("client.Recipe.Visit error : ", err)
		return nil, err
	}

	res := ParseRecipePage(doc)
	return res, nil

}

func ParseRecipePage(doc *goquery.Document) *Recipe {
	var recipe *Recipe

	// 菜谱名 链接
	name := strings.TrimSpace(doc.Find(".page-title").Text())
	link := doc.Url.String()

	// 封面
	main := doc.Find(".main-panel")
	imgNode := main.Find(".recipe-show .image img").First()
	imgSrc, exits := imgNode.Attr("src")
	if exits {
		imgSrc = strings.Split(imgSrc, "?")[0]
	}

	// x.x综合评分  xx人做过这道菜
	var score float64
	var cooked int

	scoreNode := main.Find(".stats .score .number").First()
	score, _ = strconv.ParseFloat(scoreNode.Text(), 64)

	cookedNode := main.Find(".stats .cooked .number").First()
	cooked, _ = strconv.Atoi(cookedNode.Text())

	// 获取作者信息 作者可能多个
	authors := make([]*Author, 0)
	main.Find("div.author a.avatar").Each(func(i int, s *goquery.Selection) {
		// 包含img的a元素才存在作者信息
		cloestImg := s.Find("img").First()
		if cloestImg.Size() == 0 {
			return
		}

		authorLink, exits := s.Attr("href")
		if exits {
			authorLink, _ = UrlRelativeToAbsolute(authorLink)
		}

		authorName := strings.TrimSpace(s.Text())
		avatar, _ := cloestImg.Attr("src")

		author := &Author{
			Name:   authorName,
			Link:   authorLink,
			Avatar: avatar,
		}
		authors = append(authors, author)

	})

	// 简介
	desc := strings.TrimSpace(main.Find("div.desc").First().Text())

	// 用料
	materials := make([]*RecipeMaterial, 0)
	main.Find("div.ings tr").Each(func(b int, ing *goquery.Selection) {
		materiaNamelNode := ing.Find("td.name")
		materialName := strings.TrimSpace(materiaNamelNode.Text())

		materialLink, exit := materiaNamelNode.Find("a").First().Attr("href")
		if exit {
			materialLink, _ = UrlRelativeToAbsolute(materialLink)
		}

		materiaUnitlNode := ing.Find("td.unit")
		materialUnit := strings.TrimSpace(materiaUnitlNode.Text())

		materials = append(materials, &RecipeMaterial{
			Name: materialName,
			Link: materialLink,
			Uint: materialUnit,
		})
	})

	//步骤
	steps := make([]*RecipeStep, 0)
	main.Find("div.steps li").Each(func(i int, s *goquery.Selection) {

		desc := strings.TrimSpace(s.Find("p").Text())
		img, _ := s.Find("img").Attr("src")

		step := &RecipeStep{
			Step: i + 1,
			Img:  desc,
			Desc: img,
		}
		steps = append(steps, step)
	})

	// 小贴士
	tip := strings.TrimSpace(main.Find("div.tip").Text())

	// 分类
	cates := make([]*Category, 0)
	doc.Find(".recipe-cats a").Each(func(i int, s *goquery.Selection) {
		var no, name, link string
		name = strings.TrimSpace(s.Text())
		link, exits := s.Attr("href")
		if exits {
			// 从link提取no
			no = strings.Trim(link, "/")
			p := strings.Split(no, "/")
			no = p[len(p)-1]

			link, _ = UrlRelativeToAbsolute(link)
		}
		cates = append(cates, &Category{
			No:   no,
			Name: name,
			Link: link,
		})
	})

	recipe = &Recipe{
		Name:      name,
		Link:      link,
		Cover:     imgSrc,
		Authors:   authors,
		Score:     score,
		Cooked:    cooked,
		Desc:      desc,
		Materials: materials,
		Steps:     steps,
		Tip:       tip,
		Category:  cates,
	}

	return recipe
}
