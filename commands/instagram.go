package commands

import (
	"fmt"
	"log"
	"net/http"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"

	simplejson "github.com/bitly/go-simplejson"
)

func NewTagInstagram(notifiers []NotifierType, tag string, mediaOnly bool) CommandType {
	return &baseInstagram{
		notifiers: notifiers,
		mediaOnly: mediaOnly,
		RootPath:  "tag",
		ID:        "tag-" + tag,
		PathGenerator: func() string {
			return fmt.Sprintf("https://www.instagram.com/explore/tags/%s/?__a=1", tag)
		},
	}
}

func NewUserInstagram(notifiers []NotifierType, username string, mediaOnly bool) CommandType {
	return &baseInstagram{
		notifiers: notifiers,
		mediaOnly: mediaOnly,
		RootPath:  "user",
		ID:        username,
		PathGenerator: func() string {
			return fmt.Sprintf("https://www.instagram.com/%s/?__a=1", username)
		},
	}
}

type baseInstagram struct {
	notifiers     []NotifierType
	mediaOnly     bool
	ID            string
	RootPath      string
	PathGenerator func() string
}

func (c *baseInstagram) GetName() string {
	return "instagram-" + c.ID
}

func (c *baseInstagram) GetInterval() time.Duration {
	return time.Minute * 60
}

func (c *baseInstagram) Fetch() (results []*Item, err error) {
	path := c.PathGenerator()
	json, err := getJSON(path)
	if LogIfErr(err) {
		return
	}

	data := json.GetPath(c.RootPath, "media", "nodes")

	for i := 0; i < len(data.MustArray([]interface{}{})); i++ {
		d := data.GetIndex(i)

		code, item := c.nodeToItem(d)
		if code == "" || item == nil {
			continue
		}

		items := c.tryGetMultiItemsFromNode(code, item)
		if len(items) > 0 {
			results = append(results, items...)
		} else {
			results = append(results, item)
		}

		// fmt.Println(i, item)
	}

	return
}

func (c *baseInstagram) tryGetMultiItemsFromNode(code string, item *Item) (results []*Item) {
	path := fmt.Sprintf("https://www.instagram.com/p/%s/?__a=1", code)
	json, err := getJSON(path)
	if LogIfErr(err) {
		return
	}

	data := json.GetPath("graphql", "shortcode_media", "edge_sidecar_to_children", "edges")
	if len(data.MustArray([]interface{}{})) == 0 {
		// fmt.Println("not multi")
		return
	}

	for i := 0; i < len(data.MustArray([]interface{}{})); i++ {
		d := data.GetIndex(i).GetPath("node")

		urls, id := urlsAndIDFromNode(d)
		if len(urls) == 0 || id == "" {
			continue
		}

		item = c.newItem(id, item.Desc, item.Ref, item.Created, urls)
		results = append(results, item)

		// fmt.Println(i, item)
	}

	return
}

func urlsAndIDFromNode(d *simplejson.Json) (image_urls []string, id string) {
	id = d.GetPath("id").MustString("")
	if id == "" {
		return
	}

	image_url := d.GetPath("display_src").MustString("")
	if image_url == "" {
		image_url = d.GetPath("display_url").MustString("")
	}

	if image_url == "" {
		return
	}

	image_urls = append(image_urls, image_url)

	return
}

func (c *baseInstagram) nodeToItem(d *simplejson.Json) (code string, item *Item) {
	code = d.GetPath("code").MustString("")
	if code == "" {
		return
	}

	urls, id := urlsAndIDFromNode(d)
	if len(urls) == 0 || id == "" {
		return
	}

	link := "https://www.instagram.com/p/" + code

	desc := d.GetPath("caption").MustString("")
	if desc != "" {
		desc += "\n"
	}
	desc += link

	createdUnix := d.GetPath("date").MustInt64(0)
	if createdUnix == 0 {
		return
	}

	created := time.Unix(createdUnix, 0)

	if c.mediaOnly {
		desc = ""
		created = time.Time{}
	}

	item = c.newItem(id, desc, link, created, urls)

	return
}

func getJSON(path string) (*simplejson.Json, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}

	log.Println("GET", path)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	json, err := simplejson.NewFromReader(resp.Body)
	return json, err
}

func (c *baseInstagram) newItem(id, desc, link string, created time.Time, image_urls []string) *Item {
	return &Item{
		Name:       c.GetName(),
		Identifier: "instagram_" + id,
		Desc:       desc,
		Ref:        link,
		Created:    created,
		Images:     image_urls,
		Notifiers:  c.notifiers,
	}
}
