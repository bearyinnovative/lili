package main

import (
	"log"

	. "github.com/bearyinnovative/lili/commands"
	. "github.com/bearyinnovative/lili/model"

	"github.com/bearyinnovative/lili/commands/house"
)

func main() {
	cmds := []CommandType{
		NewHackerNewsSlack(),
		NewHackerNewsAll(),
	}

	cmds = append(cmds, house.GetAllDealCommands()...)
	cmds = append(cmds, GetAllZhihuCommands()...)
	cmds = append(cmds, GetAllV2EXCommands()...)
	cmds = append(cmds, GetAllInstagramCommands()...)
	cmds = append(cmds, ArkdomeDoubanStatus)

	commander := NewCommander(cmds)
	err := commander.Run()

	if err != nil {
		log.Fatal(err)
	}
}
