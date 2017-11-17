package commands

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/util"

	"github.com/PuerkitoBio/goquery"
)

type ArkDome struct {
	notifiers []NotifierType
}

func (c *ArkDome) Name() string {
	return "arkdome-douban"
}

func (c *ArkDome) Interval() time.Duration {
	return time.Minute * 15
}

func (c *ArkDome) Notifiers() []NotifierType {
	return c.notifiers
}

func NewArkDome() *ArkDome {
	return &ArkDome{
		notifiers: CatNotifiers,
	}
}

// curl 'https://www.douban.com/people/arkdome/statuses' -H 'Accept-Encoding: gzip, deflate, br' -H 'Accept-Language: en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7,ja;q=0.6,zh-TW;q=0.5' -H 'Upgrade-Insecure-Requests: 1' -H 'User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.94 Safari/537.36' -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8' -H 'Referer: https://accounts.douban.com/login?alias=crysheen%40gmail.com&redir=https%3A%2F%2Fwww.douban.com%2Fpeople%2Farkdome%2Fstatuses&source=None&error=1013' -H 'Cookie: bid=30x_IMPMSFs; gr_user_id=76a53100-3fd1-415a-a769-167351e50e59; viewed="27016301_27085711"; _vwo_uuid_v2=7A971287F01EF3539E9B6F48E5F528F2|d19bb11fa88af1310d308ac44e76077d; ps=y; ue="crysheen@gmail.com"; dbcl2="1636924:/PQkD1bta1M"; ck=zRzn; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1510882569%2C%22https%3A%2F%2Faccounts.douban.com%2Flogin%3Falias%3Dcrysheen%2540gmail.com%26redir%3Dhttps%253A%252F%252Fwww.douban.com%252Fpeople%252Farkdome%252Fstatuses%26source%3DNone%26error%3D1013%22%5D; _pk_id.100001.8cb4=7d64afede7891844.1510882569.1.1510882569.1510882569.; _pk_ses.100001.8cb4=*; push_noty_num=0; push_doumail_num=0' -H 'Connection: keep-alive' -H 'Cache-Control: max-age=0' --compressed
func (c *ArkDome) Fetch() (results []*Item, err error) {
	// https://www.douban.com/people/arkdome/statuses (GET https://www.douban.com/people/arkdome/statuses)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "https://www.douban.com/people/arkdome/statuses", nil)

	// Headers
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
	// req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Cookie", "bid=30x_IMPMSFs; gr_user_id=76a53100-3fd1-415a-a769-167351e50e59; viewed=\"27016301_27085711\"; _vwo_uuid_v2=7A971287F01EF3539E9B6F48E5F528F2|d19bb11fa88af1310d308ac44e76077d; ps=y; ue=\"crysheen@gmail.com\"; dbcl2=\"1636924:/PQkD1bta1M\"; ck=zRzn; _pk_ref.100001.8cb4=%5B%22%22%2C%22%22%2C1510882569%2C%22https%3A%2F%2Faccounts.douban.com%2Flogin%3Falias%3Dcrysheen%2540gmail.com%26redir%3Dhttps%253A%252F%252Fwww.douban.com%252Fpeople%252Farkdome%252Fstatuses%26source%3DNone%26error%3D1013%22%5D; _pk_id.100001.8cb4=7d64afede7891844.1510882569.1.1510882569.1510882569.; _pk_ses.100001.8cb4=*; push_noty_num=0; push_doumail_num=0")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36")
	req.Header.Add("Accept-Language", "en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4,ja;q=0.2,zh-TW;q=0.2")

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if LogIfErr(err) {
		return
	}

	// Log(doc.Text())

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
			pic := s2.AttrOr("data-median-src", "")
			if pic == "" {
				pic = s2.AttrOr("src", "")
			}

			// fmt.Printf("pic %d: %s\n", i2, pic)
			if pic != "" {
				pics = append(pics, pic)
			}
		})

		// fmt.Printf("Review %d: %s[%d] (%v)\n", i, text, len(pics), created)

		if len(pics) == 0 {
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
