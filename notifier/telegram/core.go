package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type textMediaPair struct {
	text  string
	media string
}

func (t *textMediaPair) mediaType() string {
	if strings.HasSuffix(t.media, ".mp4") {
		return "video"
	}

	return "photo"
}

func toPairs(text string, media []string) []*textMediaPair {
	if len(media) == 0 {
		return []*textMediaPair{&textMediaPair{text, ""}}
	}

	results := make([]*textMediaPair, len(media))
	for i := 0; i < len(media); i++ {
		results[i] = &textMediaPair{
			text, media[i],
		}
	}

	return results
}

func (n *Notifier) notify(pairs []*textMediaPair) error {
	count := len(pairs)

	if count == 0 {
		return errors.New("nothing to send")
	}

	mediaCount := 0
	for i := 0; i < count; i++ {
		if pairs[i].media != "" {
			mediaCount += 1
		}
	}

	var method string
	var values = map[string]interface{}{}
	values["chat_id"] = n.ChatID

	// filter text only
	newPairs := []*textMediaPair{}
	for i := 0; i < count; i++ {
		if pairs[i].media == "" {
			method = "sendMessage"

			values["text"] = pairs[0].text
			if n.ParseMode != "" {
				values["parse_mode"] = n.ParseMode
			}

			values["disable_web_page_preview"] = n.DisableWebPagePreview
			values["disable_notification"] = n.DisableNotification

			err := n.send(method, values)
			if err != nil {
				log.Println("err:", err, values)
				continue
			}
		} else if pairs[i].mediaType() == "video" {
			// avoid send video error
			method = "sendVideo"
			values["video"] = pairs[i].media
			values["caption"] = pairs[i].text

			err := n.send(method, values)
			if err != nil {
				log.Println("err:", err, values)
				continue
			}
		} else {
			newPairs = append(newPairs, pairs[i])
		}
	}

	count = len(newPairs)

	if count == 0 {
		return nil
	} else if count > 1 {
		method = "sendMediaGroup"

		mediaPhotos := make([]map[string]string, count)
		for i := 0; i < count; i++ {
			mediaPhotos[i] = map[string]string{
				"type":    newPairs[i].mediaType(),
				"media":   newPairs[i].media,
				"caption": newPairs[i].text,
			}
		}

		values["media"] = mediaPhotos
	} else { // count == 1
		if newPairs[0].media == "" {
			// already handled
		} else if strings.HasSuffix(newPairs[0].media, "gif") {
			method = "sendDocument"
			values["document"] = newPairs[0].media
			values["caption"] = newPairs[0].text
		} else if strings.HasSuffix(newPairs[0].media, "mp4") {
			// already handled
		} else {
			method = "sendPhoto"
			values["photo"] = newPairs[0].media
			values["caption"] = newPairs[0].text
		}
	}

	return n.send(method, values)
}

func (n *Notifier) send(method string, values interface{}) error {
	jsonValue, err := json.Marshal(values)
	if err != nil {
		return err
	}
	body := bytes.NewBuffer(jsonValue)

	// Create client
	client := &http.Client{}

	// Create request
	path := fmt.Sprintf("https://api.telegram.org/bot%s/%s", n.Token, method)
	req, err := http.NewRequest("POST", path, body)
	if err != nil {
		return err
	}

	// Headers
	req.Header.Add("Content-Type", "application/json")

	err = req.ParseForm()
	if err != nil {
		return err
	}

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("status code error: %d %s %s", resp.StatusCode, string(b), values))
	}

	return nil
}
