package commands

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	. "../model"
	. "../util"

	"github.com/PuerkitoBio/goquery"
	simplejson "github.com/bitly/go-simplejson"
)

type BCZhihu struct {
	notifier NotifierType
}

func (c *BCZhihu) Name() string {
	return "zhihu-bearychat"
}

func (c *BCZhihu) Interval() time.Duration {
	return time.Minute * 45
}

func (c *BCZhihu) Notifier() NotifierType {
	return c.notifier
}

func NewBCZhihu() *BCZhihu {
	return &BCZhihu{
		notifier: DefaultChannelNotifier("不是真的lili"),
	}
}

func (z *BCZhihu) Fetch() (results []*Item, err error) {
	Log("start fetching", z.Name())

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.zhihu.com/r/search?q=bearychat&type=content", nil)
	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	defer resp.Body.Close()
	json, err := simplejson.NewFromReader(resp.Body)
	if LogIfErr(err) {
		return
	}

	htmls := json.GetPath("htmls")

	for i := 0; i < len(htmls.MustArray([]interface{}{})); i++ {
		h := htmls.GetIndex(i).MustString("")
		if h == "" {
			Log("can't find html in json")
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(h))
		if LogIfErr(err) {
			return nil, err
		}
		title := doc.Find(".title").Text()
		link := doc.Find("link").AttrOr("href", "")
		if link != "" {
			if !strings.HasPrefix(link, "http") {
				link = "https://www.zhihu.com" + link
			}
		} else {
			// no answer
			continue
		}

		author := doc.Find("a.author").Text()

		// use link as part of identifier
		item := z.createItem(link, fmt.Sprintf("%s\n%s: %s", title, author, link), link)
		results = append(results, item)
	}

	fmt.Printf("%v\n", results)
	return
}

func (z *BCZhihu) createItem(id, desc, ref string) *Item {
	return &Item{
		Name:       z.Name(),
		Identifier: "bc_zhihu_" + id,
		Desc:       desc,
		Ref:        ref,
	}
}
