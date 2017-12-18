package beary_commands

import (
	. "github.com/bearyinnovative/lili/notifier"
	. "github.com/bearyinnovative/lili/notifier/bearychat"
)

var LiliNotifiers []NotifierType
var CatNotifiers []NotifierType

func init() {
	LiliNotifiers = []NotifierType{
		BCChannelNotifier("不是真的lili"),
	}

	CatNotifiers = []NotifierType{
		BCChannelNotifier("云养猫"),
	}
}

func BCChannelNotifier(to string) NotifierType {
	return &IncomingNotifier{
		Domain:    "=bw52O",
		Token:     "08c0d225efc37cb33d31d089b91233d1",
		ToChannel: to,
	}
}

func BCUserNotifier(to string) NotifierType {
	return &IncomingNotifier{
		Domain: "=bw52O",
		Token:  "08c0d225efc37cb33d31d089b91233d1",
		ToUser: to,
	}
}
