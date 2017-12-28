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

type BaseV2EX struct {
	Notifiers []NotifierType
	Query     string
}

func (c *BaseV2EX) GetName() string {
	return "v2ex-" + c.Query
}

func (c *BaseV2EX) GetInterval() time.Duration {
	return time.Minute * 45
}

func (c *BaseV2EX) Fetch() (results []*Item, err error) {
	// custom search v2ex (GET https://www.googleapis.com/customsearch/v1?key=AIzaSyC1Q3F9GsEaIaxLe4zRwMeOhhNr7axtXEg&cx=011777316675351136864:22g5hinnt0i&q=bearychat)

	// Create client
	client := &http.Client{}

	// Create request
	path := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=AIzaSyC1Q3F9GsEaIaxLe4zRwMeOhhNr7axtXEg&cx=011777316675351136864:22g5hinnt0i&q=%s", url.PathEscape(c.Query))
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
			Identifier: "bc_v2ex_" + link,
			Desc:       desc,
			Ref:        link,
			Created:    created,
			Notifiers:  c.Notifiers,
		}
		results = append(results, item)
	}

	return
}
