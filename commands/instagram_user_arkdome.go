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
			notifier: CatNotifier,
			RootPath: "user",
			ID:       "arkdome",
			PathGenerator: func(token string) string {
				return fmt.Sprintf("https://www.instagram.com/%s/?__a=1", "arkdome")
			},
		},
	}
}
