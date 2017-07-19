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
			notifier: DefaultChannelNotifier("不是真的lili"),
			Query:    "Slack",
		},
	}
}
