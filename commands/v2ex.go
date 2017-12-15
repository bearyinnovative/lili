package commands

import (
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
			notifiers: LiliNotifiers,
			Query:     keyword,
		})
	}

	return
}
