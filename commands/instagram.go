package commands

import (
	"fmt"
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
	// Create client
	client := &http.Client{}

	// Create request
	path := c.PathGenerator()
	req, err := http.NewRequest("GET", path, nil)
	// fmt.Println("GET", path)

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

	data := json.GetPath(c.RootPath, "media", "nodes")

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

		if c.mediaOnly {
			desc = ""
			created = time.Time{}
		}

		item := &Item{
			Name:       c.GetName(),
			Identifier: "instagram_" + id,
			Desc:       desc,
			Ref:        link,
			Created:    created,
			Images:     image_urls,
			Notifiers:  c.notifiers,
		}
		results = append(results, item)

		// fmt.Println(i, item)
	}

	return
}
