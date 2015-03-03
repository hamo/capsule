package instance

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	"github.com/hamo/capsule/catalog"
)

type InstanceInfo struct {
	Name string `json:"name"`

	Kernel  string `json:"kernel"`
	Initrd  string `json:"initrd"`
	Cmdline string `json:"cmdline"`

	ExportConsole bool `json:"exportConsole"`

	Catalog *catalog.CatalogDir `json:"-"`
}

func New(name string) *InstanceInfo {
	return &InstanceInfo{
		Name: name,

		// FIXME:
		Kernel: "/boot/vmlinuz-linux-lts",
		Initrd: "/home/hamo/workspace/initrd.gz",

		ExportConsole: true,
	}
}

func (i *InstanceInfo) Create() error {
	// FIXME:
	qemu, err := exec.LookPath("qemu-system-x86_64")
	if err != nil {
		return errors.New("can not find qemu")
	}

	cmd := exec.Command(qemu, "-enable-kvm", "-nographic")

	if i.Kernel != "" {
		cmd.Args = append(cmd.Args, "-kernel", i.Kernel)
	}

	if i.Initrd != "" {
		cmd.Args = append(cmd.Args, "-initrd", i.Initrd)
	}

	if i.ExportConsole {
		consoleLog, err := i.Catalog.File("console.log")
		if err == nil {
			cmd.Args = append(cmd.Args, "-chardev", "file,id=console,path="+consoleLog)
			cmd.Args = append(cmd.Args, "-serial", "chardev:console")
			i.Cmdline = i.Cmdline + " " + "console=ttyS0 ignore_loglevel"
		}
	}

	controlSocket, err := i.Catalog.File("control.sock")
	if err == nil {
		cmd.Args = append(cmd.Args,
			"-device",
			"virtio-serial",
			"-chardev",
			fmt.Sprintf("socket,path=%s,server,nowait,id=control", controlSocket),
			"-device",
			"virtserialport,chardev=control,nr=1",
		)
	}

	// Cmdline
	cmd.Args = append(cmd.Args, "-append", i.Cmdline)

	// Save info
	infoFile, err := i.Catalog.File("info.json")
	if err != nil {
		return err
	}

	b, err := json.Marshal(*i)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(infoFile, b, 0600)
	if err != nil {
		return err
	}

	return cmd.Start()
}
