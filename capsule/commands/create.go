package commands

import (
	"github.com/Sirupsen/logrus"
)

func init() {
	if _, ok := CommandsList["create"]; ok {
		panic("double command register")
	}
	CommandsList["create"] = CapsuleCommand{
		Handler:     cmdCreate,
		Description: "Create a new capsule",
	}
}

func cmdCreate(args []string, logger *logrus.Logger) error {
	return nil
}
