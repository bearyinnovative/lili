package util

import (
	"fmt"
	"log"
	"os"
)

var errLogger *log.Logger

func init() {
	errLogger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

	log.Println("init util")
}

func LogIfErr(err error) bool {
	if err == nil {
		return false
	}

	errLogger.Output(2, fmt.Sprintln("[ERROR]", err))
	// panic(err)
	return true
}
