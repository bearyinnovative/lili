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
			notifier: LiliNotifier,
			Query:    "BearyChat",
		},
	}
}
