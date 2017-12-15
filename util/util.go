package util

import "log"

func init() {
	// log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("init util")
}

func LogIfErr(err error) bool {
	if err == nil {
		return false
	}

	log.Println("[ERROR]", err)
	// panic(err)
	return true
}
