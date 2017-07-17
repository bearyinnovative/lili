package model

import (
	"time"
)

type CommandType interface {
	Name() string
	Fetch() ([]*Item, error)
	Interval() time.Duration
	Notifier() NotifierType
}
