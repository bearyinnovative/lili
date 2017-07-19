package commands

import (
	. "github.com/bearyinnovative/lili/model"
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
