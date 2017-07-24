package commands

import (
	"fmt"

	. "github.com/bearyinnovative/lili/model"
)

type ArkDomeInstagram struct {
	*BaseInstagram
}

func NewArkDomeInstagram() *ArkDomeInstagram {
	return &ArkDomeInstagram{
		&BaseInstagram{
			notifiers: CatNotifiers,
			RootPath:  "user",
			ID:        "arkdome",
			PathGenerator: func(token string) string {
				return fmt.Sprintf("https://www.instagram.com/%s/?__a=1", "arkdome")
			},
		},
	}
}
