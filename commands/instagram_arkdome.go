package commands

import (
	"fmt"

	. "../model"
)

type ArkDomeInstagram struct {
	*BaseInstagram
}

func NewArkDomeInstagram() *ArkDomeInstagram {
	return &ArkDomeInstagram{
		&BaseInstagram{
			notifier: DefaultChannelNotifier("云养猫"),
			ID:       "arkdome",
			PathGenerator: func(token string) string {
				return fmt.Sprintf("https://www.instagram.com/%s/?__a=1", "arkdome")
			},
		},
	}
}
