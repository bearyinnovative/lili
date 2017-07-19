package commands

import (
	. "../model"
)

type DingDingZhihu struct {
	*BaseZhihu
}

func NewDingDingZhihu() *DingDingZhihu {
	return &DingDingZhihu{
		&BaseZhihu{
			notifier: DefaultChannelNotifier("不是真的lili"),
			Query:    "钉钉",
		},
	}
}
