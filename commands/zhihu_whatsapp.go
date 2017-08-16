package commands

import (
	. "github.com/bearyinnovative/lili/model"
)

type WhatsAppZhihu struct {
	*BaseZhihu
}

func NewWhatsAppZhihu() *WhatsAppZhihu {
	return &WhatsAppZhihu{
		&BaseZhihu{
			notifiers: LiliNotifiers,
			Query:     "WhatsApp",
		},
	}
}
