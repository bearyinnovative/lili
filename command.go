package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type Command struct {
	Name      string
	Interval  time.Duration
	ticker    *time.Ticker
	notifiers []NotifierType
}

func (c *Command) historyFileName() string {
	return "_history/" + c.Name + ".log"
}

const (
	DefaultInterval = time.Minute * 2
)

func MakeCommand(f os.FileInfo) *Command {
	interval := parseInterval(f.Name())

	return &Command{
		Name:     f.Name(),
		Interval: interval,
	}
}

// notifications.1m.py
func parseInterval(name string) time.Duration {
	comps := strings.Split(name, ".")
	if len(comps) < 3 {
		return DefaultInterval
	}

	// 1m
	timeStr := strings.ToLower(comps[1])
	timeStrLen := len(timeStr)
	if timeStrLen < 2 {
		return DefaultInterval
	}

	numStr := timeStr[0 : timeStrLen-1]
	temp, err := strconv.Atoi(numStr)
	if LogIfErr(err) {
		return DefaultInterval
	}

	num := time.Duration(temp)
	if num <= 0 {
		num = DefaultInterval
	}

	switch timeStr[timeStrLen-1:] {
	case "s":
		return num * time.Second
	case "m":
		return num * time.Minute
	case "h":
		return num * time.Hour
	case "d":
		return num * time.Hour * 24
	default:
		log.Println("can't parse", timeStr)
		return DefaultInterval
	}

	return DefaultInterval
}

func (c *Command) AddNotifier(n NotifierType) {
	c.notifiers = append(c.notifiers, n)
}

func (c *Command) Start() {
	// trigger once
	c.doJob()

	c.Stop()
	c.ticker = time.NewTicker(c.Interval)
	go func() {
		for _ = range c.ticker.C {
			c.doJob()
		}
	}()
}

func (c *Command) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
}

func (c *Command) doJob() {
	log.Printf("ticking %s: ", c.Name)
	defer fmt.Println("\n------")

	result, err := ExecBash("./" + c.Name)
	if LogIfErr(err) {
		return
	}
	hresult := c.historyResult()

	if result == hresult {
		fmt.Printf("no changes")
		return
	}

	if hresult == "" {
		log.Printf("no history")
		c.notify(result, false)
	} else {
		diff := Diff(hresult, result)
		if diff == "" {
			fmt.Printf("no changes from diff")
		} else {
			c.notify(diff, true)
		}
	}

	c.writeToHistory(result)
}

func (c *Command) notify(text string, diffed bool) {
	if diffed {
		info, err := os.Stat(c.historyFileName())
		FatalIfErr(err)

		text = fmt.Sprintf("---[M] since %v---\n%s", info.ModTime(), text)
	} else {
		text = fmt.Sprintf("---[NEW]---\n%s", text)
	}

	for i := 0; i < len(c.notifiers); i++ {
		c.notifiers[i].Notify(text)
	}
}

func (c *Command) writeToHistory(text string) {
	err := ioutil.WriteFile(c.historyFileName(), []byte(text), 0644)
	FatalIfErr(err)
}

func (c *Command) historyResult() string {
	b, err := ioutil.ReadFile(c.historyFileName())
	if err != nil {
		return ""
	}
	return string(b)
}
