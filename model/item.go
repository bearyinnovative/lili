package model

import "time"

type Item struct {
	// required
	Identifier string `bson:"identifier"`
	Name       string `bson:"name"`

	// optional
	Desc   string   `bson:"desc"`
	Ref    string   `bson:"ref"`
	Images []string `bson:"images"`

	Created time.Time `bson:"created"`
	Updated time.Time `bson:"updated"`

	// not used for now
	Key        string   `bson:"key"`
	KeyHistory []string `bson:"key_history"`
}

func (i *Item) IsValid() bool {
	if i.Identifier == "" || i.Name == "" {
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
