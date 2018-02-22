package commands

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"

	simplejson "github.com/bitly/go-simplejson"
)

type BaseGoogle struct {
	Notifiers  []NotifierType
	Identifier string
	Query      string
	Key        string
	Cx         string
}

func (c *BaseGoogle) GetName() string {
	return "google-" + c.Identifier
}

func (c *BaseGoogle) GetInterval() time.Duration {
	return time.Minute * 110
}

func (c *BaseGoogle) Fetch() (results []*Item, err error) {
	// custom search (GET https://www.googleapis.com/customsearch/v1?key=KEY_HERE&cx=CONTEXT_HERE&q=bearychat)

	// Create client
	client := &http.Client{}

	// Create request
	path := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=%s&q=%s", c.Key, c.Cx, url.PathEscape(c.Query))
	req, err := http.NewRequest("GET", path, nil)
	if LogIfErr(err) {
		return
	}

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

	htmls := json.GetPath("items")

	for i := 0; i < len(htmls.MustArray([]interface{}{})); i++ {
		hi := htmls.GetIndex(i)

		title, err := hi.GetPath("title").String()
		if LogIfErr(err) {
			continue
		}

		link, err := hi.GetPath("link").String()
		if LogIfErr(err) {
			continue
		}

		var created time.Time
		dateStr := hi.
			GetPath("pagemap").
			GetPath("metatags").
			GetIndex(0).
			GetPath("article:published_time").
			MustString("")

		if dateStr != "" {
			// Mon Jan 2 15:04:05 -0700 MST 2006
			// 2016-03-23T10:39:24Z
			created, err = time.Parse(time.RFC3339, dateStr)
			if LogIfErr(err) {
				continue
			}
		}

		desc := fmt.Sprintf("%s\n%s", title, link)
		item := &Item{
			Name: c.GetName(),
			// use link as part of identifier
			Identifier: c.GetName() + "-" + link,
			Desc:       desc,
			Ref:        link,
			Created:    created,
			Notifiers:  c.Notifiers,
		}
		results = append(results, item)
	}

	return
}
