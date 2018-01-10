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

func (n *OpenAPINotifier) Notify(id, text string, media []string) error {
	var attachment bc.MessageAttachment
	if len(media) > 0 {
		attachment.Images = []bc.MessageAttachmentImage{}
		for _, img := range media {
			attachment.Images = append(attachment.Images, bc.MessageAttachmentImage{&img})
		}
	}

	opt := &bc.MessageCreateOptions{
		VChannelID:  n.ToVChannelID,
		Text:        text,
		Attachments: []bc.MessageAttachment{attachment},
	}

	n.client.Message.Create(context.TODO(), opt)

	return nil
}
