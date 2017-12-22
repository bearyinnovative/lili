package lili

import (
	"errors"
	"log"
	"os"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/util"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dbContext DatabaseType

type DatabaseType interface {
	// return created, key_changed, error
	UpsertItem(*Item) (bool, bool, error)
	MarkNotified(*Item, bool) error
}

type Database struct {
	itemColl *mgo.Collection
}

func init() {
	mongoServer := os.Getenv("MONGO_SERVER")
	if mongoServer == "" {
		mongoServer = "localhost"
	}
	log.Println("db before dial:", mongoServer)
	session, err := mgo.Dial(mongoServer)
	if LogIfErr(err) {
		panic(err)
	}
	// defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	dbContext = &Database{
		itemColl: session.DB("lili").C("items"),
	}

	log.Println("mongo setup success")
}

func (db *Database) UpsertItem(h *Item) (bool, bool, error) {
	keyChanged := false
	if !h.IsValid() {
		return false, keyChanged, errors.New("item invalid")
	}

	query := bson.M{
		"identifier": h.Identifier,
	}
	count, err := db.itemColl.Find(query).Count()
	if LogIfErr(err) {
		return false, keyChanged, err
	}
	if count > 1 {
		return false, keyChanged, errors.New("more than one item with same identifier")
	}
	if count == 0 {
		if h.Created.IsZero() {
			h.Created = time.Now()
		}

		err = db.itemColl.Insert(h)
		if LogIfErr(err) {
			return false, keyChanged, err
		}

		return true, keyChanged, nil
	}

	var old *Item
	err = db.itemColl.Find(query).One(&old)
	if LogIfErr(err) {
		return false, keyChanged, err
	}

	h.Updated = time.Now()

	if old.Key != h.Key {
		// log.Println("key updated")
		keyChanged = true
		h.KeyHistory = append(old.KeyHistory, old.Key)
	} else {
		h.KeyHistory = old.KeyHistory
	}

	err = db.itemColl.Update(query, h)
	if LogIfErr(err) {
		return false, keyChanged, err
	}

	return false, keyChanged, nil
}

func (db *Database) MarkNotified(item *Item, notified bool) error {
	query := bson.M{
		"identifier": item.Identifier,
	}

	var t time.Time
	if notified {
		t = time.Now()
	} else {
		t = time.Time{} // empty time means haven't notified
	}
	return db.itemColl.Update(query, bson.M{"$set": bson.M{"notified_at": t}})
}
