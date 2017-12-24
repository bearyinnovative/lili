package model

import (
	"fmt"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
)

type ItemFlag int

const (
	// will override JustNotify, usually used for track history in db
	DoNotNotify ItemFlag = 1 << iota
	// normal used for force notify something and not save to db
	JustNotify ItemFlag = 1 << iota
	// won't check how old this item is created
	DoNotCheckTooOld ItemFlag = 1 << iota
)

type Item struct {
	// required
	Identifier string `bson:"identifier"`
	Name       string `bson:"name"`

	// optional
	NotifyText string   // will use this instead of desc to notify
	Desc       string   `bson:"desc"`
	Ref        string   `bson:"ref"`
	Images     []string `bson:"images"`

	Created    time.Time `bson:"created"`
	Updated    time.Time `bson:"updated"`
	NotifiedAt time.Time `bson:"notified_at"`

	Key        string   `bson:"key"`
	KeyHistory []string `bson:"key_history"`

	ItemFlags ItemFlag
}

func (i *Item) IsValid() bool {
	if i.Identifier == "" || i.Name == "" {
		return false
	}

	return true
}

func (i *Item) keyHistoryDesc() string {
	if len(i.KeyHistory) < 10 {
		return strings.Join(i.KeyHistory, "->")
	}

	results := i.KeyHistory[:2]
	results = append(results, "...")
	results = append(results, i.KeyHistory[len(i.KeyHistory)-5:]...)
	return strings.Join(results, "->")
}

func (i *Item) CheckNeedNotify(created, keyChanged bool) bool {
	if i.ItemFlags&DoNotNotify > 0 {
		return false
	}

	if i.ItemFlags&JustNotify > 0 {
		return true
	}

	checkTooOldPassed := (i.ItemFlags&DoNotCheckTooOld > 0) || i.InDays(31)

	if created || keyChanged {
		return checkTooOldPassed
	}

	return false
}

func (i *Item) GetNotifyText(created, keyChanged bool) string {
	if i.NotifyText != "" {
		return i.NotifyText
	} else if keyChanged {
		return fmt.Sprintf("%s (%s)", i.Desc, i.keyHistoryDesc())
	} else {
		return fmt.Sprintf("%s (%s)", i.Desc, humanize.Time(i.Created))
	}
}

func (i *Item) NeedSaveToDB() bool {
	// only need save to db when not just notify
	return i.ItemFlags&JustNotify == 0
}

func (i *Item) InDays(n int) bool {
	if time.Now().Sub(i.Created).Hours() < float64(n*24) {
		return true
	}

	return false
}
