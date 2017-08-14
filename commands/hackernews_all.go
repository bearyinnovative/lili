package commands

import . "github.com/bearyinnovative/lili/model"

type HackerNewsAll struct {
	*BaseHackerNews
}

func NewHackerNewsAll() *HackerNewsAll {
	return &HackerNewsAll{
		&BaseHackerNews{
			notifiers: []NotifierType{DefaultChannelNotifier("rocry_news")},
			name:      "rocry",
			shouldNotify: func(item *HNItem) bool {
				if item.Score < 50 || item.Comments < 5 {
					return false
				}

				return true
			},
		},
	}
}
