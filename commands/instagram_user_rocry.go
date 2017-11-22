package commands

import (
	"fmt"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
)

type RoCryInstagram struct {
	*BaseInstagram
}

func NewRoCryInstagram() *RoCryInstagram {
	return &RoCryInstagram{
		&BaseInstagram{
			notifiers: []NotifierType{BCUserNotifier("rocry")},
			RootPath:  "user",
			ID:        "rocry",
			PathGenerator: func(token string) string {
				return fmt.Sprintf("https://www.instagram.com/%s/?__a=1", "rocry")
			},
		},
	}
}
