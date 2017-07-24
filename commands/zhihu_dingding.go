package commands

import (
	. "github.com/bearyinnovative/lili/model"
)

type DingDingZhihu struct {
	*BaseZhihu
}

func NewDingDingZhihu() *DingDingZhihu {
	return &DingDingZhihu{
		&BaseZhihu{
			notifiers: LiliNotifiers,
			Query:     "钉钉",
		},
	}
}
