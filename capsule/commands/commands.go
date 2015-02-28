package commands

import (
	"github.com/hamo/capsule/catalog"

	"github.com/Sirupsen/logrus"
)

type CapsuleCommand struct {
	Handler     CommandHander
	Description string
}

type CommandEnv struct {
	BaseCatalog *catalog.CatalogDir

	Logger *logrus.Logger
}

type CommandHander func(args []string, cmdEnv *CommandEnv) error

var (
	// FIXME: the order of range call
	CommandsList map[string]*CapsuleCommand
)

func init() {
	CommandsList = make(map[string]*CapsuleCommand)
}

func RegisterCommand(command string, c *CapsuleCommand) {
	if _, ok := CommandsList[command]; ok {
		panic("double command register")
	}
	CommandsList[command] = c
}
