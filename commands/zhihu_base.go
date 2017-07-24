package commands

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/util"

	"github.com/PuerkitoBio/goquery"
	simplejson "github.com/bitly/go-simplejson"
)

type BaseZhihu struct {
	notifiers []NotifierType
	Query     string
}

func (c *BaseZhihu) Name() string {
	return "zhihu-" + c.Query
}

func (c *BaseZhihu) Interval() time.Duration {
	return time.Minute * 45
}

func (c *BaseZhihu) Notifiers() []NotifierType {
	return c.notifiers
}

func (z *BaseZhihu) Fetch() (results []*Item, err error) {
	client := &http.Client{}
	path := fmt.Sprintf("https://www.zhihu.com/r/search?q=%s&type=content", url.PathEscape(z.Query))
	req, err := http.NewRequest("GET", path, nil)
	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	// bytes, err := ioutil.ReadAll(resp.Body)
	// LogIfErr(err)
	// fmt.Println("testtttt", string(bytes))

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

		var created time.Time
		rawDateStr := doc.Find("a.time.text-muted").Text()
		dateStrComps := strings.Split(rawDateStr, " ")
		if len(dateStrComps) != 0 {
			loc, err := time.LoadLocation("Local")
			if LogIfErr(err) {
				return nil, err
			}

			dateStr := dateStrComps[len(dateStrComps)-1]
			created, err = time.ParseInLocation("2006-01-02", dateStr, loc)
			if LogIfErr(err) {
				return nil, err
			}
		}

		author := doc.Find("a.author").Text()

		// use link as part of identifier
		item := z.createItem(link, fmt.Sprintf("%s\n%s: %s", title, author, link), link, created)
		results = append(results, item)
	}

	return
}

func (z *BaseZhihu) createItem(id, desc, ref string, created time.Time) *Item {
	return &Item{
		Name:       z.Name(),
		Identifier: "bc_zhihu_" + id, // bc_ for history reason
		Desc:       desc,
		Ref:        ref,
		Created:    created,
	}
}
