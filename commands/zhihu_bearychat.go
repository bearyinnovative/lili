package commands

import (
	. "../model"
)

type BearyChatZhihu struct {
	*BaseZhihu
}

func NewBearyChatZhihu() *BearyChatZhihu {
	return &BearyChatZhihu{
		&BaseZhihu{
			notifier: DefaultChannelNotifier("不是真的lili"),
			Query:    "BearyChat",
		},
	}
}
