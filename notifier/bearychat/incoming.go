package bearychat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type IncomingNotifier struct {
	URL       string `yaml:"url"`
	ToUser    string `yaml:"to_user,omitempty"`
	ToChannel string `yaml:"to_channel,omitempty"`
}

/*
{
    "text": "text, this field may accept markdown",
    "markdown": true,
    "channel": "bearychat-dev",
    "attachments": [
        {
            "title": "title_1",
            "text": "attachment_text",
            "color": "#ffa500",
            "images": [
                {"url": "http://img3.douban.com/icon/ul15067564-30.jpg"}
            ]
        }
    ]
}
*/
func (n *IncomingNotifier) Notify(id, text string, media []string) error {
	path := n.URL

	dic := map[string]interface{}{
		"text": text,
	}
	if n.ToUser != "" {
		dic["user"] = n.ToUser
	}
	if n.ToChannel != "" {
		dic["channel"] = n.ToChannel
	}

	if len(media) > 0 {
		mediaArr := []interface{}{}
		for _, img := range media {
			mediaArr = append(mediaArr, map[string]string{
				"url": img,
			})
		}
		dic["attachments"] = []interface{}{
			map[string]interface{}{
				"images": mediaArr,
			},
		}
	}

	jsonValue, err := json.Marshal(dic)
	if err != nil {
		return err
	}

	body := bytes.NewBuffer(jsonValue)

	client := &http.Client{}
	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("status code error: %d", resp.StatusCode))
	}

	return nil
}
