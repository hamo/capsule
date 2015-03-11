package instance

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os/exec"
	"strconv"

	"github.com/hamo/capsule/catalog"
)

type InstanceInfo struct {
	Name string `json:"name"`

	Kernel  string `json:"kenrel"`
	Cmdline string `json:"cmdline"`

	MemorySize int `json:"memorySize"`

	ExportConsole bool `json:"exportConsole"`

	SysinitDir string `json:"-"`

	InstanceCatalog *catalog.CatalogDir `json:"-"`
	KernelCatalog   *catalog.CatalogDir `json:"-"`
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

	// Performance tuning
	// FIXME: configurable
	cmd.Args = append(cmd.Args, "-cpu", "qemu64,+ssse3,+sse4.1,+sse4.2,+x2apic")

	cmd.Args = append(cmd.Args, "-m", strconv.Itoa(i.MemorySize))

	kernelPath, err := i.KernelCatalog.TryFile("vmlinux", false)
	if err != nil {
		return errors.New("can not read vmlinux file.")
	}
	cmd.Args = append(cmd.Args, "-kernel", kernelPath)

	initrdPath, err := i.KernelCatalog.TryFile("initrd", false)
	if err == nil {
		cmd.Args = append(cmd.Args, "-initrd", initrdPath)
	}

	if i.ExportConsole {
		consoleLog, err := i.InstanceCatalog.File("console.log", false)
		if err == nil {
			cmd.Args = append(cmd.Args, "-chardev", "file,id=console,path="+consoleLog)
			cmd.Args = append(cmd.Args, "-serial", "chardev:console")
			// FIXME: ignore_loglevel?
			i.Cmdline = i.Cmdline + " " + "console=ttyS0 ignore_loglevel"
		}
	}

	controlSocket, err := i.InstanceCatalog.File("control.sock", false)
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

	// capsuled
	cmd.Args = append(cmd.Args,
		"-fsdev",
		"local,security_model=passthrough,id=sysinit,readonly,path="+i.SysinitDir,
		"-device",
		"virtio-9p-pci,fsdev=sysinit,mount_tag=sysinit",
	)

	// Cmdline
	cmd.Args = append(cmd.Args, "-append", i.Cmdline)

	// Save info
	infoFile, err := i.InstanceCatalog.File("info.json", false)
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
