package commands

import . "github.com/bearyinnovative/lili/model"

type HackerNewsSlack struct {
	*BaseHackerNews
}

func NewHackerNewsSlack() *HackerNewsSlack {
	return &HackerNewsSlack{
		&BaseHackerNews{
			notifiers: LiliNotifiers,
			keyword:   "slack",
		},
	}
}
