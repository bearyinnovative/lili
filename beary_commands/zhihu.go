package beary_commands

import (
	. "github.com/bearyinnovative/lili/commands"
	. "github.com/bearyinnovative/lili/model"
)

func GetAllZhihuCommands() (results []CommandType) {
	keywords := []string{
		"BearyChat",
		"钉钉",
		"iMessage",
		"Slack",
		"Telegram",
		"WhatsApp",
	}

	for _, keyword := range keywords {
		results = append(results, &BaseZhihu{
			Notifiers: LiliNotifiers,
			Query:     keyword,
		})
	}

	return
}
