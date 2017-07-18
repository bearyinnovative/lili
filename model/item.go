package model

import "time"

type Item struct {
	// required
	Identifier string `bson:"identifier"`
	Desc       string `bson:"desc"`
	Name       string `bson:"name"`

	Ref    string   `bson:"ref"`
	Images []string `bson:"images"`

	Key        string   `bson:"key"`
	KeyHistory []string `bson:"key_history"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`
}

func (i *Item) IsValid() bool {
	if i.Identifier == "" || i.Name == "" || i.Desc == "" {
		return false
	}

	return true
}

func (i *Item) InDays(n int) bool {
	if time.Now().Sub(i.Created).Hours() < float64(n*24) {
		return true
	}

	return false
}
