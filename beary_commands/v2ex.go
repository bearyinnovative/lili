package beary_commands

import (
	. "github.com/bearyinnovative/lili/commands"
	. "github.com/bearyinnovative/lili/model"
)

func GetAllV2EXCommands() (results []CommandType) {
	data := []string{
		"BearyChat",
		"钉钉",
		"Slack",
	}

	for _, keyword := range data {
		results = append(results, &BaseV2EX{
			Notifiers: LiliNotifiers,
			Query:     keyword,
		})
	}

	return
}
