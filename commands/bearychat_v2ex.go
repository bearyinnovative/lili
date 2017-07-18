package commands

import (
	"fmt"
	"net/http"
	"time"

	. "../model"
	. "../util"

	simplejson "github.com/bitly/go-simplejson"
)

type BCV2ex struct {
	notifier NotifierType
}

func (c *BCV2ex) Name() string {
	return "v2ex-bearychat"
}

func (c *BCV2ex) Interval() time.Duration {
	return time.Minute * 45
}

func (c *BCV2ex) Notifier() NotifierType {
	return c.notifier
}

func NewBCV2ex() *BCV2ex {
	return &BCV2ex{
		notifier: DefaultChannelNotifier("不是真的lili"),
	}
}

func (c *BCV2ex) Fetch() (results []*Item, err error) {
	// custom search v2ex (GET https://www.googleapis.com/customsearch/v1?key=AIzaSyC1Q3F9GsEaIaxLe4zRwMeOhhNr7axtXEg&cx=011777316675351136864:22g5hinnt0i&q=bearychat)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "https://www.googleapis.com/customsearch/v1?key=AIzaSyC1Q3F9GsEaIaxLe4zRwMeOhhNr7axtXEg&cx=011777316675351136864:22g5hinnt0i&q=bearychat", nil)
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
			Name: c.Name(),
			// use link as part of identifier
			Identifier: "bc_v2ex_" + link,
			Desc:       desc,
			Ref:        link,
			Created:    created,
		}
		results = append(results, item)
	}

	return
}
