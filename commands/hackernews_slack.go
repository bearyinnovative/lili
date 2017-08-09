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
				return strings.Contains(strings.ToLower(item.Title), "slack")
			},
		},
	}
}
