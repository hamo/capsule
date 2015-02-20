package commands

import (
	"github.com/Sirupsen/logrus"
)

type CapsuleCommand struct {
	Handler     CommandHander
	Description string
}

type CommandHander func(args []string, logger *logrus.Logger) error

var (
	CommandsList map[string]CapsuleCommand
)

func init() {
	CommandsList = make(map[string]CapsuleCommand)
}
