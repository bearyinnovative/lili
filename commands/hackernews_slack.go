package commands

import (
	"strings"

	. "github.com/bearyinnovative/lili/model"
)

type HackerNewsSlack struct {
	*BaseHackerNews
}

func NewHackerNewsSlack() *HackerNewsSlack {
	return &HackerNewsSlack{
		&BaseHackerNews{
			notifiers: LiliNotifiers,
			name:      "slack",
			shouldNotify: func(item *HNItem) bool {
				return checkContains(item.Title, []string{"slack", "telegram", "whatsapp"})
			},
		},
	}
}

func checkContains(title string, keywords []string) bool {
	lowerTitle := strings.ToLower(title)

	for _, key := range keywords {
		if strings.Contains(lowerTitle, key) {
			return true
		}
	}

	return false
}
