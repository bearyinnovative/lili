package bearychat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type IncomingNotifier struct {
	Domain    string
	Token     string
	ToUser    string
	ToChannel string
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
func (n *IncomingNotifier) Notify(text string, images []string) error {
	path := fmt.Sprintf("https://hook.bearychat.com/%s/incoming/%s", n.Domain, n.Token)

	dic := map[string]interface{}{
		"text": text,
	}
	if n.ToUser != "" {
		dic["user"] = n.ToUser
	}
	if n.ToChannel != "" {
		dic["channel"] = n.ToChannel
	}

	if len(images) > 0 {
		imagesArr := []interface{}{}
		for _, img := range images {
			imagesArr = append(imagesArr, map[string]string{
				"url": img,
			})
		}
		dic["attachments"] = []interface{}{
			map[string]interface{}{
				"images": imagesArr,
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
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
