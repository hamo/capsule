package commands

import (
	"flag"

	"github.com/hamo/capsule/instance"
)

var createCommand CapsuleCommand = CapsuleCommand{
	Handler:     cmdCreate,
	Description: "Create a new capsule",
}

func init() {
	RegisterCommand("create", &createCommand)
}

func cmdCreate(args []string, cmdEnv *CommandEnv) error {
	logger := cmdEnv.Logger

	createFlag := flag.NewFlagSet("create command", flag.ExitOnError)

	var flInstanceName string

	createFlag.StringVar(&flInstanceName, "name", "", "instance name")

	createFlag.Parse(args)

	if flInstanceName == "" {
		logger.Fatalln("name")
	}

	instancesCatalog, err := cmdEnv.BaseCatalog.Dir("instances")
	if err != nil {
		logger.Fatal(err)
	}

	if _, err := instancesCatalog.TryDir(flInstanceName); err == nil {
		logger.Fatalf("%s exists\n", flInstanceName)
	}

	myInstanceCatalog, err := instancesCatalog.Dir(flInstanceName)
	if err != nil {
		logger.Fatalln("create instance catalog failed.")
	}

	i := instance.New(flInstanceName)
	i.Catalog = myInstanceCatalog

	err = i.Create()
	if err != nil {
		panic(err)
	}
	return nil
}

func createUsage() {

}
