package model

import (
	"fmt"
	"strings"
	"time"

	. "github.com/bearyinnovative/lili/notifier"
	humanize "github.com/dustin/go-humanize"
)

type ItemFlag int

const (
	// normal used for force notify something and not save to db
	JustNotify ItemFlag = 1 << iota
	// won't check how old this item is created
	DoNotCheckTooOld ItemFlag = 1 << iota
)

type Item struct {
	// required
	Identifier string `bson:"identifier"`

	// optional
	Name       string   `bson:"name,omitempty"`
	NotifyText string   `bson:"notify_text,omitempty"` // will use this instead of desc to notify
	Desc       string   `bson:"desc,omitempty"`
	Ref        string   `bson:"ref,omitempty"`
	Images     []string `bson:"images,omitempty"`

	Created time.Time `bson:"created,omitempty"`
	Updated time.Time `bson:"updated,omitempty"`

	Notifiers  []NotifierType `bson:"-"`
	NotifiedAt time.Time      `bson:"notified_at,omitempty"`

	Key        string   `bson:"key,omitempty"`
	KeyHistory []string `bson:"key_history,omitempty"`

	ItemFlags ItemFlag `bson:"item_flags,omitempty"`
}

func (i *Item) IsValid() bool {
	if i.Identifier == "" {
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

func (i *Item) GetValidNotifiers(created, keyChanged bool) []NotifierType {
	if i.ItemFlags&JustNotify > 0 {
		return i.Notifiers
	}

	checkTooOldPassed := (i.ItemFlags&DoNotCheckTooOld > 0) || i.InDays(31)

	if (created || keyChanged) && checkTooOldPassed {
		return i.Notifiers
	}

	return nil
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
