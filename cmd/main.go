package main

import (
	"aws-ssh/internal/commands"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	cmds, err := commands.RegisteredCommands()
	err = cmds.Execute()
	if err != nil {
		log.Error("Error executing CLI")
		return 1
	}
	return 0
}
