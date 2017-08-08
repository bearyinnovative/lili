package commands

import . "github.com/bearyinnovative/lili/model"

type HackerNewsAll struct {
	*BaseHackerNews
}

func NewHackerNewsAll() *HackerNewsAll {
	return &HackerNewsAll{
		&BaseHackerNews{
			notifiers: []NotifierType{DefaultUserNotifier("rocry")},
			keyword:   "",
		},
	}
}
