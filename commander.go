package lili

import (
	"errors"
	"fmt"
	"log"
	"time"

	. "github.com/bearyinnovative/lili/model"
	. "github.com/bearyinnovative/lili/util"

	"github.com/dustin/go-humanize"
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
	if err != nil {
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

		if !item.InDays(31) {
			log.Println("too old to notify:", item.Desc)
			continue
		}

		notifiedCount += 1

		// notify
		text := fmt.Sprintf("%s (%s)", item.Desc, humanize.Time(item.Created))
		for _, n := range c.GetNotifiers() {
			err = n.Notify(text, item.Images)
			LogIfErr(err)
			if err == nil {
				err = DBContext.MarkNotified(item)
				LogIfErr(err)
			}
		}
	}

	log.Printf("[%s] fetched %d items, notified %d", c.GetName(), len(items), notifiedCount)
}
