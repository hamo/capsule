// FIXME: import net package so not a static binary
// FIXME: build with command:
// FIXME: CGO_ENABLED=0 go build -a -installsuffix cgo github.com/hamo/capsule/capsuled

package main

import (
	"flag"
	"fmt"
	"net/rpc"
	"os"

	"github.com/hamo/capsule/control/server"
)

var (
	flControl string
)

func init() {
	flag.StringVar(&flControl, "control", "/dev/vport0p1", "Control Port")
}

func main() {
	flag.Parse()

	fmt.Println("capsuled!")

	control, err := os.OpenFile(flControl, os.O_RDWR, 0660)
	if err != nil {
		fmt.Println(err)
	}
	defer control.Close()

	server := server.NewServer()

	rpcServer := rpc.NewServer()

	rpcServer.Register(server)

	rpcServer.ServeConn(control)

	// FIXME:
	for {
	}

	// unreachable
	panic("unreachable")
}
