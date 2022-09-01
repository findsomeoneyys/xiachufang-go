package xiachufang

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	Client    *http.Client
	UserAgent string
}

func NewClient() Client {
	return Client{
		Client:    &http.Client{},
		UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/15.4 Safari/605.1.15",
	}
}

func (c *Client) Visit(URL string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		log.Fatal("http.NewRequest error : ", err)
		return nil, err
	}

	req.Header = http.Header{
		"User-Agent": {c.UserAgent},
	}

	resp, err := c.Client.Do(req)

	if err != nil {
		log.Fatal("client.Visit error : ", err)
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("client.Visit NewDocumentFromReader error  : ", err)
		return nil, err
	}
	doc.Url = req.URL

	return doc, nil
}
