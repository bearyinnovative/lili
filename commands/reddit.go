package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

var imageExts = []string{
	"gif",
	"png",
	"jpg",
	"jepg",
	"bmp",
}

type RedditResp struct {
	Kind string `json:"kind"`
	Data struct {
		Modhash         string `json:"modhash"`
		WhitelistStatus string `json:"whitelist_status"`
		Children        []struct {
			Kind string     `json:"kind"`
			Data *ChildData `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type ChildData struct {
	Domain              string `json:"domain"`
	SubredditID         string `json:"subreddit_id"`
	ThumbnailWidth      int    `json:"thumbnail_width"`
	Subreddit           string `json:"subreddit"`
	Selftext            string `json:"selftext"`
	IsRedditMediaDomain bool   `json:"is_reddit_media_domain"`
	ID                  string `json:"id"`
	Archived            bool   `json:"archived"`
	Clicked             bool   `json:"clicked"`
	Author              string `json:"author"`
	NumCrossposts       int    `json:"num_crossposts"`
	Saved               bool   `json:"saved"`
	CanModPost          bool   `json:"can_mod_post"`
	IsCrosspostable     bool   `json:"is_crosspostable"`
	Pinned              bool   `json:"pinned"`
	Score               int    `json:"score"`
	Over18              bool   `json:"over_18"`
	Hidden              bool   `json:"hidden"`
	Preview             struct {
		Images []struct {
			Source struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"source"`
			Resolutions []struct {
				URL    string `json:"url"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"resolutions"`
			Variants struct {
			} `json:"variants"`
			ID string `json:"id"`
		} `json:"images"`
		Enabled bool `json:"enabled"`
	} `json:"preview"`
	Thumbnail             string  `json:"thumbnail"`
	ContestMode           bool    `json:"contest_mode"`
	Gilded                int     `json:"gilded"`
	Downs                 int     `json:"downs"`
	BrandSafe             bool    `json:"brand_safe"`
	PostHint              string  `json:"post_hint"`
	AuthorFlairText       string  `json:"author_flair_text"`
	Stickied              bool    `json:"stickied"`
	CanGild               bool    `json:"can_gild"`
	ThumbnailHeight       int     `json:"thumbnail_height"`
	ParentWhitelistStatus string  `json:"parent_whitelist_status"`
	Name                  string  `json:"name"`
	Spoiler               bool    `json:"spoiler"`
	Permalink             string  `json:"permalink"`
	SubredditType         string  `json:"subreddit_type"`
	Locked                bool    `json:"locked"`
	HideScore             bool    `json:"hide_score"`
	Created               float64 `json:"created"`
	URL                   string  `json:"url"`
	WhitelistStatus       string  `json:"whitelist_status"`
	Quarantine            bool    `json:"quarantine"`
	Title                 string  `json:"title"`
	CreatedUtc            float64 `json:"created_utc"`
	SubredditNamePrefixed string  `json:"subreddit_name_prefixed"`
	Ups                   int     `json:"ups"`
	NumComments           int     `json:"num_comments"`
	IsSelf                bool    `json:"is_self"`
	Visited               bool    `json:"visited"`
	IsVideo               bool    `json:"is_video"`
	Distinguished         string  `json:"distinguished"`
}

type Reddit struct {
	Subreddit string // CNY, USD, ...

	// optional
	Interval    int // in minutes
	ImageOnly   bool
	MinUpsRatio float64 // 0~1
	MinScore    int

	Notifiers []NotifierType
}

func (c *Reddit) GetName() string {
	return "reddit-" + c.Subreddit
}

func (c *Reddit) GetInterval() time.Duration {
	return time.Minute * time.Duration(c.Interval)
}

func (c *Reddit) Fetch() (results []*Item, err error) {
	path := fmt.Sprintf("https://www.reddit.com/r/%s.json", c.Subreddit)

	client := &http.Client{}

	req, err := http.NewRequest("GET", path, nil)
	if LogIfErr(err) {
		return
	}

	req.Header.Add("User-Agent", "lili")

	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	var redditResp RedditResp
	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(&redditResp)
	if LogIfErr(err) {
		return
	}

	for _, child := range redditResp.Data.Children {
		data := child.Data

		if data.Score < c.MinScore {
			continue
		}

		ratio := float64(data.Ups) / float64(data.Ups+data.Downs)
		if ratio < c.MinUpsRatio {
			continue
		}

		// try find image
		images := data.tryGetImages()
		if c.ImageOnly && len(images) == 0 {
			continue
		}

		ref := "http://reddit.com" + data.Permalink
		item := &Item{
			Name:       c.GetName(),
			Identifier: c.GetName() + "-" + data.ID,
			Desc:       fmt.Sprintf("%s %s", data.Title, ref),
			Ref:        ref,
			Created:    time.Unix(int64(data.CreatedUtc), 0),
			Notifiers:  c.Notifiers,
			Images:     images,
		}

		results = append(results, item)
	}

	return
}

func (data *ChildData) tryGetImages() []string {
	splits := strings.Split(data.URL, ".")
	if len(splits) < 2 {
		return nil
	}
	ext := strings.ToLower(splits[len(splits)-1])

	// this doesn't work for telegram
	// if ext == "gifv" {
	// 	return []string{data.URL[:len(data.URL)-1]}
	// }

	if !stringInSlice(ext, imageExts) {
		return nil
	}

	return []string{data.URL}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
