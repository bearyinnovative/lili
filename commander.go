package main

import (
	"fmt"
	"log"
	"time"

	. "./commands"
	. "./model"
	. "./util"

	"github.com/dustin/go-humanize"
)

func RunCommander() {
	cmds := []CommandType{
		NewArkDome(),
		NewBCV2ex(),
		NewBCZhihu(),
	}

	for i := 0; i < len(cmds); i++ {
		// fmt.Printf("%+v\n", cmds[i])
		start(cmds[i])
	}

	// FIXME:
	// wait forever
	select {}
}

func start(c CommandType) {
	// trigger once
	fetchAndNotify(c)

	ticker := time.NewTicker(c.Interval())
	// schedule run
	go func() {
		for _ = range ticker.C {
			fetchAndNotify(c)
		}
	}()
}

func fetchAndNotify(c CommandType) {
	items, err := c.Fetch()
	if err != nil || len(items) == 0 {
		return
	}

	notifiedCount := 0

	for _, item := range items {
		created, err := DBContext.UpsertItem(item)
		if LogIfErr(err) {
			continue
		}

		if !created {
			continue
		}

		if !item.InDays(20) {
			Log("too old to notify:", item.Desc)
			continue
		}

		notifiedCount += 1

		// notify
		text := fmt.Sprintf("[NEW] %s (%s)", item.Desc, humanize.Time(item.Created))
		c.Notifier().Notify(text, item.Images)
	}

	log.Printf("[%s] fetched %d items, notified %d", c.Name(), len(items), notifiedCount)
}
