package commands

import (
	"net/http"
	"time"

	. "../model"
	. "../util"

	simplejson "github.com/bitly/go-simplejson"
)

const (
	// deprecated
	token = "4163129.01dbb7e.8666598be3004da1b509c24bbd57336f"
)

type BaseInstagram struct {
	notifier      NotifierType
	ID            string
	PathGenerator func(string) string
}

func (c *BaseInstagram) Name() string {
	return "instagram-" + c.ID
}

func (c *BaseInstagram) Interval() time.Duration {
	return time.Minute * 60
}

func (c *BaseInstagram) Notifier() NotifierType {
	return c.notifier
}

func (c *BaseInstagram) Fetch() (results []*Item, err error) {
	// someone's recent media (GET https://api.instagram.com/v1/users/4163129/media/recent?access_token=4163129.01dbb7e.8666598be3004da1b509c24bbd57336f)

	// Create client
	client := &http.Client{}

	// Create request
	path := c.PathGenerator(token)
	req, err := http.NewRequest("GET", path, nil)

	// Fetch Request
	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	defer resp.Body.Close()
	json, err := simplejson.NewFromReader(resp.Body)
	if LogIfErr(err) {
		return
	}

	data := json.GetPath("user", "media", "nodes")

	for i := 0; i < len(data.MustArray([]interface{}{})); i++ {
		d := data.GetIndex(i)

		image_urls := []string{}

		image_url := d.GetPath("display_src").MustString("")
		if image_url == "" {
			continue
		}
		image_urls = append(image_urls, image_url)

		if len(image_urls) == 0 {
			continue
		}

		code := d.GetPath("code").MustString("")
		if code == "" {
			continue
		}
		link := "https://www.instagram.com/p/" + code

		id := d.GetPath("id").MustString("")
		if id == "" {
			continue
		}

		desc := d.GetPath("caption").MustString("")
		if desc != "" {
			desc += "\n"
		}
		desc += link

		createdUnix := d.GetPath("date").MustInt64(0)
		if createdUnix == 0 {
			continue
		}

		created := time.Unix(createdUnix, 0)

		item := &Item{
			Name:       c.Name(),
			Identifier: "instagram_" + id,
			Desc:       desc,
			Ref:        link,
			Created:    created,
			Images:     image_urls,
		}
		results = append(results, item)

		// fmt.Println(i, item)
	}

	return
}
