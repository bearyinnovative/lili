package house

import (
	"errors"
	"log"
	"os"
	"time"

	. "github.com/bearyinnovative/lili/util"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var dealCollection *mgo.Collection
var houseCollection *mgo.Collection

func init() {
	mongoServer := os.Getenv("MONGO_SERVER")
	if mongoServer == "" {
		mongoServer = "localhost"
	}
	log.Println("house before dial:", mongoServer)
	session, err := mgo.Dial(mongoServer)
	if LogIfErr(err) {
		panic(err)
	}
	// defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	dealCollection = session.DB("house").C("deals")
	houseCollection = session.DB("house").C("houses")

	log.Println("house mongo setup success")
}

func upsertDeal(d *DealItem) (bool, error) {
	query := bson.M{
		"housecode": d.HouseCode,
	}
	count, err := dealCollection.Find(query).Count()
	if LogIfErr(err) {
		return false, err
	}
	if count > 1 {
		return false, errors.New("more than one deal with same house code")
	}
	if count == 0 {
		d.FetchedAt = time.Now()
		err = dealCollection.Insert(d)
		if LogIfErr(err) {
			return false, err
		}

		return true, nil
	}

	var old *DealItem
	err = dealCollection.Find(query).One(&old)
	if LogIfErr(err) {
		return false, err
	}

	d.FetchedAt = old.FetchedAt

	err = dealCollection.Update(query, d)
	if LogIfErr(err) {
		return false, err
	}

	return false, nil
}

// return if created OR changed
func upsertHouse(h *HouseItem) (bool, error) {
	h.UpdatedAt = time.Now()

	query := bson.M{
		"housecode": h.HouseCode,
	}
	count, err := houseCollection.Find(query).Count()
	if LogIfErr(err) {
		return false, err
	}
	if count > 1 {
		return false, errors.New("more than one deal with same house code")
	}

	if count == 0 {
		h.FetchedAt = h.UpdatedAt

		err = houseCollection.Insert(h)
		if LogIfErr(err) {
			return false, err
		}

		return true, nil
	}

	var old *HouseItem
	err = houseCollection.Find(query).One(&old)
	if LogIfErr(err) {
		return false, err
	}

	h.FetchedAt = old.FetchedAt

	// price changed, append to price history
	priceChanged := h.Price != old.Price
	if priceChanged {
		h.HistoryPrices = append(old.HistoryPrices, old.Price)
	}

	err = houseCollection.Update(query, h)
	if LogIfErr(err) {
		return false, err
	}

	return priceChanged, nil
}
