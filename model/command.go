package model

import (
	"time"

	. "github.com/bearyinnovative/lili/notifier"
)

type CommandType interface {
	GetName() string
	Fetch() ([]*Item, error)
	GetInterval() time.Duration
	GetNotifiers() []NotifierType
}
