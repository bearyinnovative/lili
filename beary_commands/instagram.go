package beary_commands

import (
	. "github.com/bearyinnovative/lili/commands"
	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/notifier"
)

func GetAllInstagramCommands() []CommandType {
	return []CommandType{
		NewTagInstagram(CatNotifiers, "arkdome"),

		NewUserInstagram(CatNotifiers, "arkdome"),
		NewUserInstagram(CatNotifiers, "dabie.cat"),
		NewUserInstagram(CatNotifiers, "matsumotoooooo"),
		NewUserInstagram([]NotifierType{BCUserNotifier("rocry")}, "rocry"),
	}
}
