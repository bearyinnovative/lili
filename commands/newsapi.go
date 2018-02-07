package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type NewsAPIResp struct {
	Status       string     `json:"status"`
	TotalResults int        `json:"totalResults"`
	Articles     []*Article `json:"articles"`
}

type Article struct {
	Source struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
}

type NewsAPISubscriber struct {
	Name         string
	Notifiers    []NotifierType
	ShouldNotify func(*Article) bool
}

type NewsAPI struct {
	Subscribers []*NewsAPISubscriber

	Endpoint string
	Sources  string
	APIKey   string
	Interval int // in minutes
}

func (c *NewsAPI) GetName() string {
	return fmt.Sprintf("newsapi-%s-%s", c.Endpoint, c.Sources)
}

func (c *NewsAPI) GetInterval() time.Duration {
	return time.Minute * time.Duration(c.Interval)
}

func (c *NewsAPI) Fetch() (results []*Item, err error) {
	// Create client
	client := &http.Client{}

	// Create request
	path := fmt.Sprintf("https://newsapi.org/v2/%s?sources=%s&apiKey=%s", c.Endpoint, c.Sources, c.APIKey)
	log.Println(path)
	req, err := http.NewRequest("GET", path, nil)
	if LogIfErr(err) {
		return
	}

	// Fetch Request
	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	var result *NewsAPIResp

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(&result)
	if LogIfErr(err) {
		return
	}

	if result.Status != "ok" {
		log.Println("status not ok:", result)
		return
	}

	for _, article := range result.Articles {
		if LogIfErr(err) {
			continue
		}

		var images []string

		// remove image for now
		// if article.URLToImage != "" {
		// 	images = []string{article.URLToImage}
		// }

		for _, sub := range c.Subscribers {
			if !sub.ShouldNotify(article) {
				continue
			}

			name := c.GetName() + "-" + sub.Name
			item := &Item{
				Name:       name,
				Identifier: name + "-" + article.URL,
				Desc:       fmt.Sprintf("[%s](%s)", article.Title, article.URL),
				Ref:        article.URL,
				Created:    article.PublishedAt,
				Images:     images,
				Notifiers:  sub.Notifiers,
			}
			results = append(results, item)
		}
	}

	return
}
