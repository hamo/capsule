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

	Kernel  string `json:"kenrel"`
	Cmdline string `json:"cmdline"`

	ExportConsole bool `json:"exportConsole"`

	Catalog       *catalog.CatalogDir `json:"-"`
	KernelCatalog *catalog.CatalogDir `josn:"-"`
}

func New(name string) *InstanceInfo {
	return &InstanceInfo{
		Name: name,

		// FIXME:
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

	kernelPath, err := i.KernelCatalog.TryFile("vmlinux", false)
	if err != nil {
		return errors.New("can not read vmlinux file.")
	}
	cmd.Args = append(cmd.Args, "-kernel", kernelPath)

	initrdPath, err := i.KernelCatalog.TryFile("initrd", false)
	if err == nil {
		cmd.Args = append(cmd.Args, "-initrd", initrdPath)

		// FIXME: initrd is too big so we need at least 512MB memory
		cmd.Args = append(cmd.Args, "-m", "512")
	}

	if i.ExportConsole {
		consoleLog, err := i.Catalog.File("console.log", false)
		if err == nil {
			cmd.Args = append(cmd.Args, "-chardev", "file,id=console,path="+consoleLog)
			cmd.Args = append(cmd.Args, "-serial", "chardev:console")
			// FIXME: ignore_loglevel?
			i.Cmdline = i.Cmdline + " " + "console=ttyS0 ignore_loglevel"
		}
	}

	controlSocket, err := i.Catalog.File("control.sock", false)
	if err == nil {
		cmd.Args = append(cmd.Args,
			"-device",
			"virtio-serial",
			"-chardev",
			"socket,server,nowait,id=control,path="+controlSocket,
			"-device",
			"virtserialport,chardev=control,nr=1",
		)
	}

	// Cmdline
	cmd.Args = append(cmd.Args, "-append", i.Cmdline)

	// Save info
	infoFile, err := i.Catalog.File("info.json", false)
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
