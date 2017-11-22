package model

import (
	"time"

	. "github.com/bearyinnovative/lili/notifier"
)

type CommandType interface {
	Name() string
	Fetch() ([]*Item, error)
	Interval() time.Duration
	Notifiers() []NotifierType
}
