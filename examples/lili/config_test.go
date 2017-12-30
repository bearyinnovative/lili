package main

import (
	"testing"
)

func TestConfigInit(t *testing.T) {
	config, err := NewConfigFromFile("./config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(config.Douban[0].ID)
	t.Log(config.V2EX[0].Notifiers[0].ToChannel)
}

func TestConfigToCommands(t *testing.T) {
	config, err := NewConfigFromFile("./config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	commands := config.ToCommandTypes()
	if len(commands) == 0 {
		t.Fatal("no commands")
	}

	t.Logf("generated %d commands\n", len(commands))

	for _, c := range commands {
		t.Logf("%s\n", c.GetName())
	}
}
