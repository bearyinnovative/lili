package commands

import (
	. "../model"
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
