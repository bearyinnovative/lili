package commands

import (
	"fmt"
	"log"
	"strings"
	"time"

	. "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type Twitter struct {
	client *Client

	Name      string
	Interval  int
	MediaOnly bool

	Query    string
	Username string

	MinFavCount int
	MinRetCount int

	// optional
	Notifiers []NotifierType
}

// when query and username both exist, will fetch by username and filter by query
func NewTwitter(name,
	consumerKey, consumerSecret, token, tokenSecret,
	query, username string, mediaOnly bool,
	minFavCount, minRetCount int,
	interval int, notifiers []NotifierType) *Twitter {

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	oauthToken := oauth1.NewToken(token, tokenSecret)
	httpClient := config.Client(oauth1.NoContext, oauthToken)

	// Twitter client
	client := NewClient(httpClient)

	return &Twitter{
		client:      client,
		Name:        name,
		Interval:    interval,
		MediaOnly:   mediaOnly,
		Query:       query,
		Username:    username,
		MinFavCount: minFavCount,
		MinRetCount: minRetCount,
		Notifiers:   notifiers,
	}
}

func (c *Twitter) GetName() string {
	return "Twitter-" + c.Name
}

func (c *Twitter) GetInterval() time.Duration {
	return time.Minute * time.Duration(c.Interval)
}

func (c *Twitter) Fetch() (results []*Item, err error) {
	var tweets []Tweet

	if c.Username != "" {
		params := &UserTimelineParams{ScreenName: c.Username}
		tweets, _, err = c.client.Timelines.UserTimeline(params)
		log.Printf("[%s] fetched %d with username: %s\n", c.GetName(), len(tweets), c.Username)

		if LogIfErr(err) {
			return
		}
	} else if c.Query != "" {
		params := &SearchTweetParams{
			Query: c.Query,
		}

		search, _, err1 := c.client.Search.Tweets(params)
		if LogIfErr(err1) {
			err = err1
			return
		}

		tweets = search.Statuses
		log.Printf("[%s] fetched %d with query: %s\n", c.GetName(), len(tweets), c.Query)
	}

	for _, t := range tweets {
		if c.Username != "" {
			if t.User.ScreenName != c.Username {
				continue
			}
		}

		if c.Query != "" {
			if !strings.Contains(t.Text, c.Query) {
				continue
			}
		}

		item := c.toItem(t)
		if item != nil {
			// log.Printf("[%s] %s %s\n", c.GetName(), t.Text, item.Ref)
			results = append(results, item)
		}
	}

	return
}

func (c *Twitter) toItem(t Tweet) *Item {
	if t.FavoriteCount < c.MinFavCount {
		log.Println("fav count:", t.FavoriteCount)
		return nil
	}

	if t.RetweetCount < c.MinRetCount {
		log.Println("ret count:", t.RetweetCount)
		return nil
	}

	var media []string = nil

	if t.Entities != nil {
		for _, m := range t.Entities.Media {

			if m.MediaURLHttps != "" {
				media = append(media, m.MediaURLHttps)
			}
		}
	}

	if c.MediaOnly && len(media) == 0 {
		return nil
	}

	ref := fmt.Sprintf("https://twitter.com/%s/status/%d", t.User.ScreenName, t.ID)

	desc := ""
	if !c.MediaOnly {
		desc = fmt.Sprintf("%s %s", t.Text, ref)
	}

	created, _ := t.CreatedAtTime()

	return &Item{
		Name:       c.GetName(),
		Identifier: c.GetName() + "-" + t.IDStr,
		Desc:       desc,
		Ref:        ref,
		Notifiers:  c.Notifiers,
		Created:    created,
		Images:     media,
	}
}
