package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/hamo/capsule/capsule/commands"
	"github.com/hamo/capsule/catalog"

	"github.com/Sirupsen/logrus"
)

var (
	flHelp  bool
	flDebug bool
)

func init() {
	flag.Usage = usage
	flag.BoolVar(&flHelp, "help", false, "usage")
	flag.BoolVar(&flDebug, "debug", false, "debug")
}

func usage() {
	fmt.Printf("Usage: capsule <args> [command] <command args>\n")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Printf("%15s     %s\n", "command", "description")
	for k, v := range commands.CommandsList {
		fmt.Printf("%15s     %s\n", k, v.Description)
	}
}

func main() {
	flag.Parse()

	if flHelp {
		usage()
		return
	}

	logger := logrus.New()
	if flDebug {
		logger.Level = logrus.DebugLevel
	} else {
		// FIXME: custom by arg
		logger.Level = logrus.InfoLevel
	}

	if runtime.GOARCH != "amd64" {
		logger.WithFields(logrus.Fields{
			"module": "main",
			"arch":   runtime.GOARCH,
		}).Fatalln("currently capsule only supports amd64.")
	}

	cmdline := flag.Args()
	if len(cmdline) == 0 {
		usage()
		return
	}

	capsuleCmd, ok := commands.CommandsList[cmdline[0]]
	if !ok {
		usage()
		return
	}

	catalogBase := os.Getenv("CAPSULE_ROOT")
	if catalogBase == "" {
		catalogBase = os.Getenv("HOME")
	}

	cmdEnv := new(commands.CommandEnv)

	baseCatalog, err := catalog.NewBaseCatalogDir(catalogBase, ".capsule")
	if err != nil {
		logger.Fatalln(err)
	}

	cmdEnv.BaseCatalog = baseCatalog
	cmdEnv.Logger = logger

	handler := capsuleCmd.Handler
	handler(cmdline[1:], cmdEnv)
}
