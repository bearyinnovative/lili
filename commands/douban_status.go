package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/util"
)

type DoubanStatus struct {
	ID        string
	Notifiers []NotifierType
}

func (c *DoubanStatus) GetName() string {
	return "douban-status-" + c.ID
}

func (c *DoubanStatus) GetInterval() time.Duration {
	return time.Minute * 15
}

func (c *DoubanStatus) Fetch() (results []*Item, err error) {
	// https://frodo.douban.com/api/v2/status/user_timeline/144859503 (GET https://frodo.douban.com/api/v2/status/user_timeline/144859503?count=15&os_rom=miui6&apikey=0dad551ec0f84ed02907ff5c42e8ec70&channel=Google_Market&udid=8f7b52865761deac6d547c8d415ed0a079704517&_sig=xA%2F56W6u7Yca1iIgMkXS3NO6Y9A%3D&_ts=1510883310)

	// Create client
	client := &http.Client{}

	// Create request
	path := fmt.Sprintf("https://frodo.douban.com/api/v2/status/user_timeline/%s?count=15&os_rom=miui6&apikey=0dad551ec0f84ed02907ff5c42e8ec70&channel=Google_Market&udid=8f7b52865761deac6d547c8d415ed0a079704517", c.ID)
	req, err := http.NewRequest("GET", path, nil)
	if LogIfErr(err) {
		return
	}

	// Headers
	req.Header.Add("Host", "frodo.douban.com")
	req.Header.Add("Connection", "Keep-Alive")
	req.Header.Add("User-Agent", "api-client/1 com.douban.frodo/5.11.0(114) Android/25 product/sagit vendor/Xiaomi model/MI 6  rom/miui6  network/wifi")

	err = req.ParseForm()
	if LogIfErr(err) {
		return
	}

	// Fetch Request
	resp, err := client.Do(req)
	if LogIfErr(err) {
		return
	}

	// // Read Response Body
	// respBody, _ := ioutil.ReadAll(resp.Body)

	// // Display Results
	// fmt.Println("response Status : ", resp.Status)
	// fmt.Println("response Headers : ", resp.Header)
	// fmt.Println("response Body : ", string(respBody))

	var result *StatusResult

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	err = decoder.Decode(&result)
	if LogIfErr(err) {
		return
	}

	loc, err := time.LoadLocation("Local")
	if LogIfErr(err) {
		return
	}

	for _, it := range result.Items {
		text := it.Status.Text
		realIt := it.Status
		if it.Status.ResharedStatus != nil {
			realIt = it.Status.ResharedStatus
			text = fmt.Sprintf("转播: %s\n%s: %s", text, realIt.Author.Name, realIt.Text)
		}

		pics := []string{}
		for _, image := range realIt.Images {
			if image.Large.URL != "" {
				pics = append(pics, image.Large.URL)
			}
		}

		created, err := time.ParseInLocation("2006-01-02 15:04:05", realIt.CreateTime, loc)
		if LogIfErr(err) {
			continue
		}

		item := &Item{
			Name:       c.GetName(),
			Identifier: c.GetName() + "-" + realIt.ID,
			Desc:       fmt.Sprintf("%s [Link](%s)", text, realIt.SharingURL),
			Ref:        realIt.SharingURL,
			Created:    created,
			Images:     pics,
			Notifiers:  c.Notifiers,
		}
		results = append(results, item)
	}

	return
}

type StatusResult struct {
	Count int           `json:"count"`
	Items []*StatusItem `json:"items"`
}

type StatusItem struct {
	Status   *Status `json:"status"`
	Type     string  `json:"type"`
	Comments []struct {
		Author struct {
			URL    string `json:"url"`
			URI    string `json:"uri"`
			Avatar string `json:"avatar"`
			ID     string `json:"id"`
			UID    string `json:"uid"`
			Kind   string `json:"kind"`
			Type   string `json:"type"`
			Loc    struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				UID  string `json:"uid"`
			} `json:"loc"`
			Name string `json:"name"`
		} `json:"author"`
		ID         int           `json:"id"`
		Text       string        `json:"text"`
		Entities   []interface{} `json:"entities"`
		CreateTime string        `json:"create_time"`
		RefComment struct {
		} `json:"ref_comment"`
		URI string `json:"uri"`
	} `json:"comments"`
}

type Status struct {
	Liked                   bool           `json:"liked"`
	SubscriptionText        string         `json:"subscription_text"`
	ForbidReshareAndComment interface{}    `json:"forbid_reshare_and_comment"`
	CreateTime              string         `json:"create_time"`
	Card                    interface{}    `json:"card"`
	Entities                []interface{}  `json:"entities"`
	ParentStatus            interface{}    `json:"parent_status"`
	LikeCount               int            `json:"like_count"`
	CommentsCount           int            `json:"comments_count"`
	Images                  []*StatusImage `json:"images"`
	IsStatusAd              bool           `json:"is_status_ad"`
	ReshareID               string         `json:"reshare_id"`
	Activity                string         `json:"activity"`
	URI                     string         `json:"uri"`
	SharingURL              string         `json:"sharing_url"`
	ID                      string         `json:"id"`
	ResharedStatus          *Status        `json:"reshared_status"`
	IsSubscription          bool           `json:"is_subscription"`
	ResharesCount           int            `json:"reshares_count"`
	Text                    string         `json:"text"`
	Author                  struct {
		URL        string `json:"url"`
		URI        string `json:"uri"`
		Avatar     string `json:"avatar"`
		ID         string `json:"id"`
		UID        string `json:"uid"`
		VerifyType int    `json:"verify_type"`
		Kind       string `json:"kind"`
		Type       string `json:"type"`
		Loc        struct {
			ID   string `json:"id"`
			Name string `json:"name"`
			UID  string `json:"uid"`
		} `json:"loc"`
		Name string `json:"name"`
	} `json:"author"`
	ResharersCount int `json:"resharers_count"`
}

type StatusImage struct {
	Large struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"large"`
	IsAnimated bool `json:"is_animated"`
	Normal     struct {
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"normal"`
}
