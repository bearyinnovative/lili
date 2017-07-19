package commands

import (
	"fmt"

	. "../model"
)

type DabieCatInstagram struct {
	*BaseInstagram
}

func NewDabieCatInstagram() *DabieCatInstagram {
	return &DabieCatInstagram{
		&BaseInstagram{
			notifier: CatNotifier,
			RootPath: "user",
			ID:       "dabie.cat",
			PathGenerator: func(token string) string {
				return fmt.Sprintf("https://www.instagram.com/%s/?__a=1", "dabie.cat")
			},
		},
	}
}
