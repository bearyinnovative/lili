package beary_commands

import (
	. "github.com/bearyinnovative/lili/commands"
	. "github.com/bearyinnovative/lili/model"
)

var ArkdomeDoubanStatus CommandType = &DoubanStatus{
	"144859503",
	CatNotifiers,
}
