package bearychat

import (
	"context"

	bc "github.com/bearyinnovative/bearychat-go/openapi"
)

type OpenAPINotifier struct {
	client       *bc.Client
	ToVChannelID string
}

func NewOpenAPINotifier(token, vid string) (*OpenAPINotifier, error) {
	client := bc.NewClient(token)

	return &OpenAPINotifier{
		client,
		vid,
	}, nil
}

func (n *OpenAPINotifier) Notify(text string, images []string) error {
	opt := &bc.MessageCreateOptions{
		VChannelID:  n.ToVChannelID,
		Text:        text,
		Attachments: nil,
	}

	n.client.Message.Create(context.TODO(), opt)

	return nil
}
