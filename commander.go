package main

import (
	"errors"
	"log"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/util"
)

type Commander struct {
	cmds []CommandType
}

func NewCommander(cmds []CommandType) *Commander {
	return &Commander{
		cmds: cmds,
	}
}

func (c *Commander) Run() error {
	if len(c.cmds) == 0 {
		return errors.New("no commands")
	}

	for i := 0; i < len(c.cmds); i++ {
		// fmt.Printf("%+v\n", cmds[i])
		start(c.cmds[i])
	}

	// FIXME:
	// wait forever
	select {}
}

func start(c CommandType) {
	// trigger once
	go func() {
		fetchAndNotify(c)
	}()

	ticker := time.NewTicker(c.GetInterval())
	// schedule run
	go func() {
		for _ = range ticker.C {
			fetchAndNotify(c)
		}
	}()
}

func fetchAndNotify(c CommandType) {
	items, err := c.Fetch()
	if LogIfErr(err) {
		return
	}

	notifiedCount := 0

	for _, item := range items {
		created, keyChanged := false, false
		var err error
		if item.NeedSaveToDB() {
			created, keyChanged, err = dbContext.UpsertItem(item)

			if LogIfErr(err) {
				continue
			}
		}

		if !item.CheckNeedNotify(created, keyChanged) {
			continue
		}

		notifiedCount += 1

		// notify text
		text := item.GetNotifyText(created, keyChanged)

		// notify
		for _, n := range c.GetNotifiers() {
			err = n.Notify(text, item.Images)
			LogIfErr(err)

			// TODO: need handle multi notifiers
			if item.NeedSaveToDB() {
				err = dbContext.MarkNotified(item, err == nil)
				LogIfErr(err)
			}
		}
	}

	log.Printf("[%s] fetched %d items, notified %d", c.GetName(), len(items), notifiedCount)
}
