package commands

import (
	"flag"
	"fmt"
	"net/rpc"
	"os"
	"path/filepath"
	"time"

	"github.com/hamo/capsule/instance"

	"github.com/Sirupsen/logrus"
)

var (
	GlobalCapsuledDir  = filepath.Join("usr", "lib", "capsule", "libexec")
	GlobalCapsuledPath = filepath.Join(GlobalCapsuledDir, "capsuled")
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

		flMemorySize int

		flCmdline string
	)

	createFlag.StringVar(&flInstanceName, "name", "", "instance name")
	createFlag.StringVar(&flKernelName, "kernel", "", "kernel name")

	createFlag.IntVar(&flMemorySize, "m", 512, "memory size")

	createFlag.StringVar(&flCmdline, "cmdline", "", "cmdline")

	createFlag.Parse(args)

	if flInstanceName == "" {
		fmt.Fprintf(os.Stderr, "Please provide new instance name.\n\n")
		createFlag.PrintDefaults()
		os.Exit(1)
	}

	if flKernelName == "" {
		fmt.Fprintf(os.Stderr, "Please provide kernel name used by new instance.\n\n")
		createFlag.PrintDefaults()
		os.Exit(1)
	}

	instancesCatalog, err := cmdEnv.BaseCatalog.Dir("instances")
	if err != nil {
		logger.Fatalf("create instances catalog failed: %s", err)
	}

	if _, err := instancesCatalog.TryDir(flInstanceName); err == nil {
		logger.Fatalf("target instance name %s exists", flInstanceName)
	}

	myInstanceCatalog, err := instancesCatalog.Dir(flInstanceName)
	if err != nil {
		logger.Fatalf("create target instance catalog failed: %s", err)
	}

	kernelsCatalog, err := cmdEnv.BaseCatalog.Dir("kernels")
	if err != nil {
		myInstanceCatalog.Cleanup(true)
		logger.Fatalf("read kernels catalog failed: %s", err)
	}

	myKernelCatalog, err := kernelsCatalog.TryDir(flKernelName)
	if err != nil {
		myInstanceCatalog.Cleanup(true)
		logger.Fatalf("read kernel %s catalog failed: %s", flKernelName, err)
	}

	i := instance.New(flInstanceName)
	i.Kernel = flKernelName
	i.Cmdline = flCmdline

	if flMemorySize < 512 {
		// FIXME: support huge initrd
		cmdEnv.Logger.Infoln("workaround: force memory size to 512M")
		i.MemorySize = 512
	} else {
		i.MemorySize = flMemorySize
	}

	i.InstanceCatalog = myInstanceCatalog
	i.KernelCatalog = myKernelCatalog

	// Find capsuled
	// it should be either placed $CAPSULE_ROOT/capsuled/capsuled or
	// /usr/lib/capsule/libexec/capusled
	// Pass its parent dir to qemu
	capsuledDir := filepath.Join(cmdEnv.BaseCatalog.Path, "capsuled")
	if fi, err := os.Stat(filepath.Join(capsuledDir, "capsuled")); err != nil || !fi.Mode().IsRegular() {
		// Can not find capsuled at $CAPSULE_ROOT/capsuled/capsuled
		if fi, err := os.Stat(GlobalCapsuledPath); err != nil || !fi.Mode().IsRegular() {
			myInstanceCatalog.Cleanup(true)
			cmdEnv.Logger.Fatalln("can not find capsuled")
		} else {
			capsuledDir = GlobalCapsuledDir
		}
	}
	i.SysinitDir = capsuledDir

	err = i.Create()
	if err != nil {
		myInstanceCatalog.Cleanup(true)
		cmdEnv.Logger.Fatalf("create instance failed: %s", err)
	}

	retryTimes := 5
	retryTicker := time.NewTicker(2 * time.Second)
	retryPass := false

	var rpcClient *rpc.Client

	for _ = range retryTicker.C {
		if retryTimes == 0 {
			break
		}

		controlSock, err := myInstanceCatalog.TryFile("control.sock", false)
		if err != nil {
			retryTimes -= 1
			cmdEnv.Logger.Debugf("try to open control socket failed: %s", err)
			continue
		}

		rpcClient, err = rpc.Dial("unix", controlSock)
		if err != nil {
			retryTimes -= 1
			cmdEnv.Logger.Debugf("try to dial control socket failed: %s", err)
			continue
		}
		cmdEnv.Logger.Debugf("dial control socket success")
		retryPass = true
		break
	}

	retryTicker.Stop()
	if !retryPass {
		cmdEnv.Logger.Debugln("leave instance catalog for debugging.")
		if cmdEnv.Logger.Level != logrus.DebugLevel {
			myInstanceCatalog.Cleanup(true)
		}
		cmdEnv.Logger.Fatalf("instance starts but can not be connected.")
	}

	var d time.Duration
	err = rpcClient.Call("Server.Alive", struct{}{}, &d)
	if err != nil {
		cmdEnv.Logger.Fatalf("rpc call failed: %s", err)
	}

	cmdEnv.Logger.Infoln("Yoo, we are alive.")

	myInstanceCatalog.Cleanup(false)

	return nil
}

func createUsage() {

}
