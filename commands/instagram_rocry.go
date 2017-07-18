package commands

import (
	"fmt"

	. "../model"
)

type RoCryInstagram struct {
	*BaseInstagram
}

func NewRoCryInstagram() *RoCryInstagram {
	return &RoCryInstagram{
		&BaseInstagram{
			notifier: DefaultUserNotifier("rocry"),
			ID:       "rocry",
			PathGenerator: func(token string) string {
				return fmt.Sprintf("https://www.instagram.com/%s/?__a=1", "rocry")
			},
		},
	}
}
