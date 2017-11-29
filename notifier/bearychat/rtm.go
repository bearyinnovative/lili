package bearychat

import (
	"log"

	bc "github.com/bearyinnovative/bearychat-go"
)

type RTMNotifier struct {
	context      *bc.RTMContext
	ToVChannelID string
}

func NewRTMNotifier(token, vid string) (*RTMNotifier, error) {
	context, err := bc.NewRTMContext(token)
	if err != nil {
		return nil, err
	}

	err, _, _ = context.Run()
	if err != nil {
		return nil, err
	}

	return &RTMNotifier{
		context,
		vid,
	}, nil
}

func (n *RTMNotifier) Notify(text string, images []string) error {
	dic := map[string]interface{}{
		"text": text,
	}

	dic["vchannel_id"] = n.ToVChannelID
	dic["type"] = "message"

	// TODO: this doesn't work for now
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

	log.Println(dic)

	err := n.context.Loop.Send(dic)
	if err != nil {
		return err
	}

	return nil
}
