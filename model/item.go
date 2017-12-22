package model

import (
	"strings"
	"time"
)

type Item struct {
	// required
	Identifier string `bson:"identifier"`
	Name       string `bson:"name"`

	// optional
	Desc   string   `bson:"desc"`
	Ref    string   `bson:"ref"`
	Images []string `bson:"images"`

	Created    time.Time `bson:"created"`
	Updated    time.Time `bson:"updated"`
	NotifiedAt time.Time `bson:"notified_at"`

	Key        string   `bson:"key"`
	KeyHistory []string `bson:"key_history"`

	DoNotCheckTooOld bool
}

func (i *Item) IsValid() bool {
	if i.Identifier == "" || i.Name == "" {
		return false
	}

	return true
}

func (i *Item) KeyHistoryDesc() string {
	if len(i.KeyHistory) < 10 {
		return strings.Join(i.KeyHistory, "->")
	}

	results := i.KeyHistory[:2]
	results = append(results, "...")
	results = append(results, i.KeyHistory[len(i.KeyHistory)-5:]...)
	return strings.Join(results, "->")
}

func (i *Item) InDays(n int) bool {
	if time.Now().Sub(i.Created).Hours() < float64(n*24) {
		return true
	}

	return false
}
