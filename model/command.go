package model

import (
	"time"
)

type CommandType interface {
	GetName() string
	Fetch() ([]*Item, error)
	GetInterval() time.Duration
}
