package beary_commands

import (
	. "github.com/bearyinnovative/lili/commands"
	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
)

func NewHackerNewsAll() CommandType {
	return &BaseHackerNews{
		Notifiers: []NotifierType{BCChannelNotifier("rocry_news")},
		Name:      "rocry",
		ShouldNotify: func(item *HNItem) bool {
			if item.Score < 50 || item.Comments < 5 {
				return false
			}

			return true
		},
	}
}
