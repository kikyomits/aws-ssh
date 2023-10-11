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
	if err != nil {
		log.Error(err.Error())
		return 1
	}
	err = cmds.Execute()
	if err != nil {
		log.Error(err.Error())
		return 1
	}
	return 0
}
