package twitter

import (
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

type Notifier struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

func (n *Notifier) Notify(text string, images []string) error {
	config := oauth1.NewConfig(n.ConsumerKey, n.ConsumerSecret)
	token := oauth1.NewToken(n.AccessToken, n.AccessSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	// Twitter client
	client := twitter.NewClient(httpClient)

	// Send a Tweet
	if len(images) > 0 {
		text += "\n" + strings.Join(images, "\n")
	}
	_, _, err := client.Statuses.Update(text, nil)
	return err
}
