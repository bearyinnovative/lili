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

	err = runAndTryKeepConnected(context)
	if err != nil {
		return nil, err
	}

	return &RTMNotifier{
		context,
		vid,
	}, nil
}

func runAndTryKeepConnected(context *bc.RTMContext) error {
	err, _, errC := context.Run()
	if err != nil {
		return err
	}

	go tryKeepConnected(context, errC)
	return nil
}

func tryKeepConnected(context *bc.RTMContext, errC chan error) {
Loop:
	for {
		select {
		case err := <-errC:
			log.Printf("rtm loop error: %+v", err)
			if err := context.Loop.Stop(); err != nil {
				log.Fatal(err)
			}

			runAndTryKeepConnected(context)
			break Loop
		}
	}
}

func (n *RTMNotifier) Notify(text string, images []string) error {
	dic := map[string]interface{}{
		"text": text,
	}

	dic["vchannel_id"] = n.ToVChannelID
	dic["type"] = "message"

	// TODO: this doesn't work for now
	/*
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
	*/

	// log.Println(dic)

	err := n.context.Loop.Send(dic)
	if err != nil {
		return err
	}

	return nil
}
