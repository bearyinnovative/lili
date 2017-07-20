package commands

import (
	"fmt"

	. "github.com/bearyinnovative/lili/model"
)

type MatsumotooooooInstagram struct {
	*BaseInstagram
}

func NewMatsumotooooooInstagram() *MatsumotooooooInstagram {
	return &MatsumotooooooInstagram{
		&BaseInstagram{
			notifier: CatNotifier,
			RootPath: "user",
			ID:       "matsumotoooooo",
			PathGenerator: func(token string) string {
				return fmt.Sprintf("https://www.instagram.com/%s/?__a=1", "matsumotoooooo")
			},
		},
	}
}
