package xiachufang

type RecipeListResult struct {
	Recipes  []*Recipe `json:"recipes"`
	Current  int       `json:"current"`
	PrevPage string    `json:"prev_page"`
	NextPage string    `json:"next_page"`
}

type Recipe struct {
	Name           string            `json:"name"`
	Link           string            `json:"link"`
	Cover          string            `json:"cover"`
	Authors        []*Author         `json:"author"`
	Score          float64           `json:"score"`
	Cooked         int               `json:"cooked"`                      // 正常搜索菜谱展示
	CookedLast7Day int               `json:"cooked_last_7_day,omitempty"` // 分类搜索时展示
	Desc           string            `json:"desc"`
	Materials      []*RecipeMaterial `json:"materials"`
	Steps          []*RecipeStep     `json:"steps"`
	Tip            string            `json:"tip"`
	Category       []*Category       `json:"category,omitempty"`
}

type Author struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Link   string `json:"link"`
}

type RecipeMaterial struct {
	Name string `json:"name"`
	Uint string `json:"unit"`
	Link string `json:"link"`
}

type RecipeStep struct {
	Step int    `json:"step"`
	Img  string `json:"img"`
	Desc string `json:"desc"`
}

type Category struct {
	No       string      `json:"no"`
	Name     string      `json:"name"`
	Link     string      `json:"link"`
	Children []*Category `json:"children,omitempty"`
}
