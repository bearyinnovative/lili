package commands

import (
	. "github.com/bearyinnovative/lili/model"
)

type IMessageZhihu struct {
	*BaseZhihu
}

func NewIMessageZhihu() *IMessageZhihu {
	return &IMessageZhihu{
		&BaseZhihu{
			notifiers: LiliNotifiers,
			Query:     "iMessage",
		},
	}
}
