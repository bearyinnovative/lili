package commands

import (
	"fmt"
	"net/url"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type Flickr struct {
	client *FlickrClient
	method string // flickr.photos.getContactsPhotos

	Name     string
	Interval int

	// extra params
	UserID string

	// optional
	Notifiers []NotifierType
}

func NewFlickr(name, method,
	consumerKey, consumerSecret, token, tokenSecret string,
	userID string,
	interval int, notifiers []NotifierType) *Flickr {
	client := NewFlickrClient(
		consumerKey,
		consumerSecret,
		token,
		tokenSecret,
	)

	return &Flickr{
		client:   client,
		method:   method,
		Name:     name,
		Interval: interval,

		UserID: userID,

		Notifiers: notifiers,
	}
}

func (c *Flickr) GetName() string {
	return "Flickr-" + c.Name
}

func (c *Flickr) GetInterval() time.Duration {
	return time.Minute * time.Duration(c.Interval)
}

func (c *Flickr) Fetch() (results []*Item, err error) {

	params := url.Values{}
	if c.UserID != "" {
		params.Set("user_id", c.UserID)
	}
	resp, err := c.client.GetWithParams(c.method, params)
	if LogIfErr(err) {
		return
	}

	for _, post := range resp.Photos.Photo {
		if post.Media != "photo" {
			continue
		}
		if post.URLL == "" {
			continue
		}

		ref := fmt.Sprintf("https://www.flickr.com/photos/%s/%s/in/feed", post.Owner, post.ID)

		item := &Item{
			Name:       c.GetName(),
			Identifier: c.GetName() + "-" + post.ID,
			Desc:       "",
			Ref:        ref,
			Notifiers:  c.Notifiers,
			Images:     []string{post.URLL},
		}

		results = append(results, item)
	}

	return
}
