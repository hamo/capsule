package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/hamo/capsule/capsule/commands"

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
	fmt.Printf("%8s     %s\n", "command", "description")
	for k, v := range commands.CommandsList {
		fmt.Printf("%8s     %s\n", k, v.Description)
	}
}

func main() {
	flag.Parse()

	if runtime.GOARCH != "amd64" {
		logrus.WithField("arch", runtime.GOARCH).Fatalln("currently capsule only supports amd64.")
	}

}
