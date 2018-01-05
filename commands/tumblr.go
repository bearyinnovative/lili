package commands

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"

	tumblr_go "github.com/tumblr/tumblrclient.go"
)

type TumblrResp struct {
	Meta struct {
		Status int    `json:"status"`
		Msg    string `json:"msg"`
	} `json:"meta"`
	Response struct {
		Posts []*TumblrPost `json:"posts"`
	} `json:"response"`
}

type TumblrPost struct {
	Type               string        `json:"type"`
	BlogName           string        `json:"blog_name"`
	ID                 int64         `json:"id"`
	PostURL            string        `json:"post_url"`
	Slug               string        `json:"slug"`
	Date               string        `json:"date"`
	Timestamp          int           `json:"timestamp"`
	State              string        `json:"state"`
	Format             string        `json:"format"`
	ReblogKey          string        `json:"reblog_key"`
	Tags               []interface{} `json:"tags"`
	ShortURL           string        `json:"short_url"`
	Summary            string        `json:"summary"`
	IsBlocksPostFormat bool          `json:"is_blocks_post_format"`
	RecommendedSource  interface{}   `json:"recommended_source"`
	RecommendedColor   interface{}   `json:"recommended_color"`
	Followed           bool          `json:"followed"`
	Liked              bool          `json:"liked"`
	NoteCount          int           `json:"note_count"`
	Caption            string        `json:"caption"`
	Reblog             struct {
		Comment  string `json:"comment"`
		TreeHTML string `json:"tree_html"`
	} `json:"reblog"`
	Trail          []interface{} `json:"trail"`
	ImagePermalink string        `json:"image_permalink,omitempty"`
	Photos         []struct {
		Caption      string `json:"caption"`
		OriginalSize struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"original_size"`
		AltSizes []struct {
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"alt_sizes"`
	} `json:"photos,omitempty"`
	CanLike          bool   `json:"can_like"`
	CanReblog        bool   `json:"can_reblog"`
	CanSendInMessage bool   `json:"can_send_in_message"`
	CanReply         bool   `json:"can_reply"`
	DisplayAvatar    bool   `json:"display_avatar"`
	PhotosetLayout   string `json:"photoset_layout,omitempty"`
	VideoURL         string `json:"video_url,omitempty"`
	HTML5Capable     bool   `json:"html5_capable,omitempty"`
	ThumbnailURL     string `json:"thumbnail_url,omitempty"`
	ThumbnailWidth   int    `json:"thumbnail_width,omitempty"`
	ThumbnailHeight  int    `json:"thumbnail_height,omitempty"`
	Duration         int    `json:"duration,omitempty"`
	Player           []struct {
		Width     int    `json:"width"`
		EmbedCode string `json:"embed_code"`
	} `json:"player,omitempty"`
	VideoType      string `json:"video_type,omitempty"`
	SourceURL      string `json:"source_url,omitempty"`
	SourceTitle    string `json:"source_title,omitempty"`
	Artist         string `json:"artist,omitempty"`
	TrackName      string `json:"track_name,omitempty"`
	AlbumArt       string `json:"album_art,omitempty"`
	Embed          string `json:"embed,omitempty"`
	Plays          int    `json:"plays,omitempty"`
	AudioURL       string `json:"audio_url,omitempty"`
	AudioSourceURL string `json:"audio_source_url,omitempty"`
	AudioType      string `json:"audio_type,omitempty"`
}

type Tumblr struct {
	client *tumblr_go.Client
	path   string // from like, dashboard

	Name      string
	Interval  int
	MediaOnly bool

	// optional
	Notifiers []NotifierType
}

func NewTumblr(name, typeName, consumerKey, consumerSecret, token, tokenSecret string, interval int, mediaOnly bool, notifiers []NotifierType) (*Tumblr, error) {
	path, err := getPathFromType(typeName)
	if LogIfErr(err) {
		return nil, err
	}

	client := tumblr_go.NewClientWithToken(
		consumerKey,
		consumerSecret,
		token,
		tokenSecret,
	)

	t := &Tumblr{
		client:    client,
		path:      path,
		Name:      name,
		Interval:  interval,
		MediaOnly: mediaOnly,
		Notifiers: notifiers,
	}

	return t, err
}

func getPathFromType(typeName string) (string, error) {
	switch typeName {
	case "dashboard":
		return "/user/dashboard", nil
	case "likes":
		return "/user/likes", nil
	default:
		return "", errors.New("can't find type: " + typeName)
	}
}

func (c *Tumblr) GetName() string {
	return "Tumblr-" + c.Name
}

func (c *Tumblr) GetInterval() time.Duration {
	return time.Minute * time.Duration(c.Interval)
}

func (c *Tumblr) Fetch() (results []*Item, err error) {
	// path := fmt.Sprintf("https://api.tumblr.com/v2%s", c.path)

	resp, err := c.client.Get(c.path)
	if LogIfErr(err) {
		return
	}

	var tumblrResp TumblrResp
	err = json.Unmarshal(resp.GetBody(), &tumblrResp)
	if LogIfErr(err) {
		return
	}

	for _, post := range tumblrResp.Response.Posts {
		media := post.getMedia()
		if media == nil {
			continue
		}

		var desc string
		var created time.Time
		if !c.MediaOnly {
			desc = post.Summary
			created = time.Unix(int64(post.Timestamp), 0)
		}

		item := &Item{
			Name:       c.GetName(),
			Identifier: c.GetName() + "-" + strconv.FormatInt(post.ID, 10),
			Desc:       desc,
			Ref:        post.PostURL,
			Created:    created,
			Notifiers:  c.Notifiers,
			Images:     media,
		}

		results = append(results, item)
	}

	return
}

func (p *TumblrPost) getMedia() []string {
	switch p.Type {
	case "video":
		return []string{p.VideoURL}
	case "photo":
		if len(p.Photos) == 0 {
			return nil
		}

		results := make([]string, len(p.Photos))
		for i, photo := range p.Photos {
			results[i] = photo.OriginalSize.URL
		}
		return results

	default:
		return nil
	}
}
