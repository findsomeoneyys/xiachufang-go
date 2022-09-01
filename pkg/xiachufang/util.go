package xiachufang

import (
	"net/url"
	"path"
)

const API_ROOT = "https://www.xiachufang.com/"

var urlMap = map[string]string{
	"SEARCH":   "/search/",
	"RECIPE":   "/recipe/",
	"CATEGORY": "/category/",
}

func GetApiURL(s string) (u *url.URL, exits bool) {
	u, _ = url.Parse(API_ROOT)
	sub, exits := urlMap[s]
	if exits {
		u.Path = path.Join(u.Path, sub)
		return
	}
	return
}

func UrlRelativeToAbsolute(rawURL string) (resURL string, err error) {
	base, _ := url.Parse(API_ROOT)

	raw, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	resURL = base.ResolveReference(raw).String()
	return resURL, nil

}
