package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type BaseHackerNews struct {
	notifiers    []NotifierType
	name         string
	shouldNotify func(*HNItem) bool
}

func (c *BaseHackerNews) Name() string {
	return "hackernews-" + c.name
}

func (c *BaseHackerNews) Interval() time.Duration {
	return time.Minute * 15
}

func (c *BaseHackerNews) Notifiers() []NotifierType {
	return c.notifiers
}

/*
{
  "by" : "andreasley",
  "descendants" : 13,
  "id" : 14955693,
  "kids" : [ 14956399, 14957207, 14956021, 14956308, 14956147, 14956151 ],
  "score" : 58,
  "time" : 1502180665,
  "title" : "A Systematic Analysis of the Juniper Dual EC Incident [pdf]",
  "type" : "story",
  "url" : "https://www.cs.uic.edu/~s/papers/juniper2016/juniper2016.pdf"
}
*/
type HNItem struct {
	Comments int    `json:"descendants"`
	ID       int    `json:"id"`
	Score    int    `json:"score"`
	Time     int    `json:"time"`
	Title    string `json:"title"`
	URL      string `json:"url"`
}

func (c *BaseHackerNews) Fetch() (results []*Item, err error) {
	// someone's recent media (GET https://api.instagram.com/v1/users/4163129/media/recent?access_token=4163129.01dbb7e.8666598be3004da1b509c24bbd57336f)

	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("GET", "https://hacker-news.firebaseio.com/v0/topstories.json", nil)

	// Fetch Request
	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)

	// bytes, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(bytes))

	var ids []int
	err = decoder.Decode(&ids)
	if LogIfErr(err) {
		return
	}

	ch := make(chan *Item)
	wg := new(sync.WaitGroup)

	// 30 for first page
	for idx, id := range ids[:30] {
		wg.Add(1)
		go func(idx, id int) {
			item := c.getItem(client, idx, id)
			if item != nil {
				ch <- item
			}
			wg.Done()
		}(idx, id)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	for item := range ch {
		results = append(results, item)
	}

	return
}

func (c *BaseHackerNews) getItem(client *http.Client, idx, id int) *Item {
	// Create request
	path := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", id)
	req, err := http.NewRequest("GET", path, nil)
	// fmt.Println("GET", path)

	// Fetch Request
	resp, err := client.Do(req)
	if LogIfErr(err) {
		return nil
	}

	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	hnItem := HNItem{}
	err = decoder.Decode(&hnItem)
	if LogIfErr(err) {
		return nil
	}

	if !c.shouldNotify(&hnItem) {
		return nil
	}

	if hnItem.Comments == 0 && hnItem.URL == "" {
		return nil
	}

	commentPath := fmt.Sprintf("https://news.ycombinator.com/item?id=%d", hnItem.ID)
	desc := fmt.Sprintf("[%s](%s)\nrank: %d, %d/[%d](%s)", hnItem.Title, hnItem.URL, idx, hnItem.Score, hnItem.Comments, commentPath)
	return &Item{
		Name:       c.Name(),
		Identifier: fmt.Sprintf("hn_%s_%d", c.name, hnItem.ID),
		Desc:       desc,
		Ref:        hnItem.URL,
		Created:    time.Unix(int64(hnItem.Time), 0),
	}
}
