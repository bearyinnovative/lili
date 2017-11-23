package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Notifier struct {
	Token  string
	ChatID int
}

func (n *Notifier) Notify(text string, images []string) error {
	imageCount := len(images)

	if imageCount > 1 {
		mediaPhotos := make([]map[string]string, imageCount)
		for i := 0; i < imageCount; i++ {
			mediaPhotos[i] = map[string]string{
				"type":    "photo",
				"media":   images[i],
				"caption": text,
			}
		}
		return n.send("sendMediaGroup", map[string]interface{}{
			"caption": text,
			"chat_id": n.ChatID,
			"media":   mediaPhotos,
		})
	} else if imageCount == 1 {
		return n.send("sendPhoto", map[string]interface{}{
			"caption": text,
			"chat_id": n.ChatID,
			"photo":   images[0],
		})
	} else {
		return n.send("sendMessage", map[string]interface{}{
			"text":    text,
			"chat_id": n.ChatID,
			// "parse_mode": "Markdown",
		})
	}
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
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
