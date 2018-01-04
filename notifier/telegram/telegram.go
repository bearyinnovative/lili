package telegram

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Notifier struct {
	Token string `yaml:"token"`
	// `@channel_name` or integer id as string: `-123456`
	ChatID    string `yaml:"chat_id"`
	ParseMode string `yaml:"parse_mode,omitempty"`
}

func (n *Notifier) Notify(text string, images []string) error {
	imageCount := len(images)

	var method string
	var values = map[string]interface{}{}
	values["chat_id"] = n.ChatID

	if imageCount > 1 {
		method = "sendMediaGroup"

		mediaPhotos := make([]map[string]string, imageCount)
		for i := 0; i < imageCount; i++ {
			mediaPhotos[i] = map[string]string{
				"type":    "photo",
				"media":   images[i],
				"caption": text,
			}
		}

		values["caption"] = text
		values["media"] = mediaPhotos
	} else if imageCount == 1 {
		if strings.HasSuffix(images[0], "gif") {
			method = "sendDocument"
			values["document"] = images[0]
		} else {
			method = "sendPhoto"
			values["photo"] = images[0]
		}

		values["caption"] = text
	} else {
		method = "sendMessage"

		values["text"] = text
		if n.ParseMode != "" {
			values["parse_mode"] = n.ParseMode
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
		return errors.New(fmt.Sprintf("status code error: %d %s", resp.StatusCode, string(b)))
	}

	return nil
}
