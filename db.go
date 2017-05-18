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
	// return created, error
	UpsertItem(*Item) (bool, error)
	MarkNotified(*Item) error
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

func (db *Database) UpsertItem(h *Item) (bool, error) {
	if !h.IsValid() {
		return false, errors.New("item invalid")
	}

	query := bson.M{
		"identifier": h.Identifier,
	}
	count, err := db.itemColl.Find(query).Count()
	if LogIfErr(err) {
		return false, err
	}
	if count > 1 {
		return false, errors.New("more than one item with same identifier")
	}
	if count == 0 {
		if h.Created.IsZero() {
			h.Created = time.Now()
		}

		err = db.itemColl.Insert(h)
		if LogIfErr(err) {
			return false, err
		}

		// h.JustCreated = true
		return true, nil
	}

	var old *Item
	err = db.itemColl.Find(query).One(&old)
	if LogIfErr(err) {
		return false, err
	}

	h.Updated = time.Now()

	if old.Key != h.Key {
		// log.Println("key updated")
		// h.KeyChanged = true
		h.KeyHistory = append(old.KeyHistory, old.Key)
	} else {
		h.KeyHistory = old.KeyHistory
	}

	err = db.itemColl.Update(query, h)
	if LogIfErr(err) {
		return false, err
	}

	return false, nil
}

func (db *Database) MarkNotified(item *Item) error {
	query := bson.M{
		"identifier": item.Identifier,
	}

	return db.itemColl.Update(query, bson.M{"$set": bson.M{"notified_at": time.Now()}})
}
