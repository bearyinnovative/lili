package commands

import (
	"fmt"

	. "github.com/bearyinnovative/lili/model"
)

type ArkDomeInstagram2 struct {
	*BaseInstagram
}

func NewArkDomeInstagram2() *ArkDomeInstagram2 {
	return &ArkDomeInstagram2{
		&BaseInstagram{
			notifier: CatNotifier,
			RootPath: "tag",
			ID:       "tag-arkdome",
			PathGenerator: func(token string) string {
				return fmt.Sprintf("https://www.instagram.com/explore/tags/%s/?__a=1", "arkdome")
			},
		},
	}
}
