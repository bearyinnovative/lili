package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

func RunCommander() {
	// enter the commands folder
	if err := os.Chdir("./commands"); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cmds := findCommands(".")
	for i := 0; i < len(cmds); i++ {
		// fmt.Printf("%+v\n", cmds[i])
		cmds[i].Start()
	}

	time.Sleep(time.Hour * 100)
}

func findCommands(path string) []*Command {
	n := &BCIncommingNotifier{
		Domain: "=bw52O",
		Token:  "08c0d225efc37cb33d31d089b91233d1",
	}

	cmds := []*Command{}
	files, err := ioutil.ReadDir(path)
	FatalIfErr(err)

	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".") {
			continue
		}
		if strings.HasPrefix(f.Name(), "_") {
			continue
		}

		cmd := MakeCommand(f)
		if cmd == nil {
			log.Println("can't make command from:", f)
			continue
		}
		cmd.AddNotifier(n)
		cmds = append(cmds, cmd)
	}

	return cmds
}
