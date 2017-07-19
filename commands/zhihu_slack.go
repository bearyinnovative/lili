package commands

import (
	. "github.com/bearyinnovative/lili/model"
)

type SlackZhihu struct {
	*BaseZhihu
}

func NewSlackZhihu() *SlackZhihu {
	return &SlackZhihu{
		&BaseZhihu{
			notifier: LiliNotifier,
			Query:    "Slack",
		},
	}
}
