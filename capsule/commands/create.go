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

	var (
		flInstanceName string
		flKernelName   string

		flCmdline string
	)

	createFlag.StringVar(&flInstanceName, "name", "", "instance name")
	createFlag.StringVar(&flKernelName, "kernel", "", "kernel name")

	createFlag.StringVar(&flCmdline, "cmdline", "", "cmdline")

	createFlag.Parse(args)

	if flInstanceName == "" {
		logger.Fatalln("name")
	}
	if flKernelName == "" {
		logger.Fatalln("kernel")
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

	kernelsCatalog, err := cmdEnv.BaseCatalog.Dir("kernels")
	if err != nil {
		logger.Fatalln("can not read kernel catalog")
	}

	myKernelCatalog, err := kernelsCatalog.TryDir(flKernelName)
	if err != nil {
		logger.Fatalf("can not find kernel named %s.", flKernelName)
	}

	i := instance.New(flInstanceName)
	i.Kernel = flKernelName
	i.Cmdline = flCmdline

	i.Catalog = myInstanceCatalog
	i.KernelCatalog = myKernelCatalog

	err = i.Create()
	if err != nil {
		panic(err)
	}
	return nil
}

func createUsage() {

}
