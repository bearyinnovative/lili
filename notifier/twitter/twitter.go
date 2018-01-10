package twitter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Notifier struct {
	ConsumerKey    string `yaml:"consumer_key"`
	ConsumerSecret string `yaml:"consumer_secret"`
	AccessToken    string `yaml:"access_token"`
	AccessSecret   string `yaml:"access_secret"`
}

func (n *Notifier) Notify(id, text string, media []string) error {
	config := oauth1.NewConfig(n.ConsumerKey, n.ConsumerSecret)
	token := oauth1.NewToken(n.AccessToken, n.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	if len(media) > 0 {
		text += "\n" + strings.Join(media, "\n")
	}
	_, resp, err := client.Statuses.Update(text, nil)

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return errors.New(fmt.Sprintf("status code error: %d", resp.StatusCode))
	}

	return err
}
