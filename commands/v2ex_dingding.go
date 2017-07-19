package commands

import (
	. "../model"
)

type DingDingV2EX struct {
	*BaseV2EX
}

func NewDingDingV2EX() *DingDingV2EX {
	return &DingDingV2EX{
		&BaseV2EX{
			notifier: LiliNotifier,
			Query:    "钉钉",
		},
	}
}
