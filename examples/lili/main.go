package main

import (
	"log"
	"os"

	. "github.com/bearyinnovative/lili"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./config.yaml"
	}
	config, err := NewConfigFromFile(configPath)
	if err != nil {
		log.Fatal(err)
	}

	commander := NewCommander(config.ToCommandTypes())
	err = commander.Run()

	if err != nil {
		log.Fatal(err)
	}
}
