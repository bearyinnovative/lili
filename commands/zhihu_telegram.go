package commands

import (
	. "github.com/bearyinnovative/lili/model"
)

type TelegramZhihu struct {
	*BaseZhihu
}

func NewTelegramZhihu() *TelegramZhihu {
	return &TelegramZhihu{
		&BaseZhihu{
			notifiers: LiliNotifiers,
			Query:     "Telegram",
		},
	}
}
