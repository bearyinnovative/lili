package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	. "../util"
)

var LiliNotifier NotifierType
var CatNotifier NotifierType

func init() {
	LiliNotifier = DefaultChannelNotifier("不是真的lili")
	CatNotifier = DefaultChannelNotifier("云养猫")
}

type NotifierType interface {
	Notify(text string, images []string)
}

type BCIncommingNotifier struct {
	Domain    string
	Token     string
	ToUser    string
	ToChannel string
}

func DefaultChannelNotifier(to string) NotifierType {
	return &BCIncommingNotifier{
		Domain:    "=bw52O",
		Token:     "08c0d225efc37cb33d31d089b91233d1",
		ToChannel: to,
	}
}

func DefaultUserNotifier(to string) NotifierType {
	return &BCIncommingNotifier{
		Domain: "=bw52O",
		Token:  "08c0d225efc37cb33d31d089b91233d1",
		ToUser: to,
	}
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
func (bc *BCIncommingNotifier) Notify(text string, images []string) {
	path := fmt.Sprintf("https://hook.bearychat.com/%s/incoming/%s", bc.Domain, bc.Token)

	dic := map[string]interface{}{
		"text": text,
	}
	if bc.ToUser != "" {
		dic["user"] = bc.ToUser
	}
	if bc.ToChannel != "" {
		dic["channel"] = bc.ToChannel
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
	// Log("jsonValue:", string(jsonValue))
	FatalIfErr(err)

	body := bytes.NewBuffer(jsonValue)

	client := &http.Client{}
	req, err := http.NewRequest("POST", path, body)
	FatalIfErr(err)

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	_, err = client.Do(req)
	LogIfErr(err)
}
