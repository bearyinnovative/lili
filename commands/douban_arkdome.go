package commands

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	. "../model"
	. "../util"

	"github.com/PuerkitoBio/goquery"
)

type ArkDome struct {
	notifier NotifierType
}

func (c *ArkDome) Name() string {
	return "arkdome-douban"
}

func (c *ArkDome) Interval() time.Duration {
	return time.Minute * 15
}

func (c *ArkDome) Notifier() NotifierType {
	return c.notifier
}

func NewArkDome() *ArkDome {
	return &ArkDome{
		notifier: CatNotifier,
	}
}

func (c *ArkDome) Fetch() (results []*Item, err error) {
	// https://www.douban.com/people/arkdome/statuses (GET https://www.douban.com/people/arkdome/statuses)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "https://www.douban.com/people/arkdome/statuses", nil)

	// Headers
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	// req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Cookie", "ll=\"118282\"; bid=tbE4t34jGnk; gr_user_id=8cf9a078-5fdd-4401-8b3a-5f50df6360f5; __yadk_uid=FdrsrjGz9MKcSKrZ8e6gLbCbQtj5lUev; __ads_session=Pije5XRM6whAD1IZSAA=; ap=1; _ga=GA1.2.753293480.1479571382; viewed=\"1223823_25932288_27019102_27019086_26638586_26699339_26931905_25964764_26663629_1049136\"; _vwo_uuid_v2=846BC4497BC86865F1F51AAD5337DEE0|fbf2b9b40cc8980c14ced40abc6aa0c9; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1500361019%2C%22https%3A%2F%2Fwww.google.com.hk%2F%22%5D; ps=y; ue=\"crysheen@gmail.com\"; dbcl2=\"1636924:w6LsCLZmE9c\"; ck=s8Nz; _pk_id.100001.8cb4=c02ab624005a54bb.1481428558.32.1500361748.1500340831.; _pk_ses.100001.8cb4=*; push_noty_num=0; push_doumail_num=0; __utmt=1; __utma=30149280.753293480.1479571382.1500343904.1500361019.89; __utmb=30149280.11.7.1500361748055; __utmc=30149280; __utmz=30149280.1500361019.89.14.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); __utmv=30149280.163")
	req.Header.Add("Referer", "https://accounts.douban.com/login?alias=crysheen%40gmail.com&redir=https%3A%2F%2Fwww.douban.com%2Fpeople%2Farkdome%2Fstatuses&source=None&error=1013")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Cache-Control", "max-age=0")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.115 Safari/537.36")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4,ja;q=0.2")

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if LogIfErr(err) {
		return
	}

	doc.Find(".status-item").Each(func(i int, s *goquery.Selection) {
		loc, err := time.LoadLocation("Local")
		if LogIfErr(err) {
			return
		}

		dateStr := s.Find(".actions span.created_at").AttrOr("title", "")
		if dateStr == "" {
			return
		}
		created, err := time.ParseInLocation("2006-01-02 15:04:05", dateStr, loc)

		link := s.Find("div.hd").AttrOr("data-status-url", "")
		if link == "" {
			return
		}

		text := s.Find(".status-saying blockquote p").Text()
		text = strings.TrimSpace(text)
		if text == "" {
			text = s.Find(".text blockquote").Text()
		}
		text = strings.TrimSpace(text) // text possible ""

		pics := []string{}
		s.Find("span.group-pic").Find("img").Each(func(i2 int, s2 *goquery.Selection) {
			pic := s2.AttrOr("src", "")
			// fmt.Printf("pic %d: %s\n", i2, pic)
			if pic != "" {
				pics = append(pics, pic)
			}
		})

		// fmt.Printf("Review %d: %s[%d] (%v)\n", i, text, len(pics), created)

		if len(pics) == 0 || text == "" {
			return
		}

		item := &Item{
			Name: c.Name(),
			// use link as part of identifier
			Identifier: "ad_douban_" + link,
			Desc:       fmt.Sprintf("%s\n%s", text, link),
			Ref:        link,
			Created:    created,
			Images:     pics,
		}
		results = append(results, item)
	})

	// Log(doc)

	return
}
