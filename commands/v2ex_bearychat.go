package commands

import (
	. "github.com/bearyinnovative/lili/model"
)

type BearyChatV2EX struct {
	*BaseV2EX
}

func NewBearyChatV2EX() *BearyChatV2EX {
	return &BearyChatV2EX{
		&BaseV2EX{
			notifier: LiliNotifier,
			Query:    "BearyChat",
		},
	}
}
