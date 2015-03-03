package commands

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"

	"github.com/hamo/capsule/kernel"
	"github.com/hamo/capsule/util"

	"github.com/Sirupsen/logrus"
)

var (
	kernelListCommand = CapsuleCommand{
		Handler:     cmdKernelList,
		Description: "List kernels",
	}

	kernelImportCommand = CapsuleCommand{
		Handler:     cmdKernelImport,
		Description: "Import kernels",
	}
	kernelExportCommand = CapsuleCommand{
		Handler:     cmdKernelExport,
		Description: "Export kernels",
	}
)

func init() {
	RegisterCommand("kernel-list", &kernelListCommand)
	RegisterCommand("kernel-import", &kernelImportCommand)
	RegisterCommand("kernel-export", &kernelExportCommand)

}

func cmdKernelExport(args []string, cmdEnv *CommandEnv) error {
	return nil
}

func cmdKernelImport(args []string, cmdEnv *CommandEnv) error {
	kernelImportFlag := flag.NewFlagSet("kernel-import command", flag.ExitOnError)

	var (
		flName    string
		flVersion string

		flVmlinux string
		flModule  string

		flKPack string
	)

	kernelImportFlag.StringVar(&flName, "name", "", "name")
	kernelImportFlag.StringVar(&flVersion, "version", "", "version")

	kernelImportFlag.StringVar(&flVmlinux, "vmlinux", "", "vmlinux path")
	kernelImportFlag.StringVar(&flModule, "module", "", "module dir path")

	kernelImportFlag.StringVar(&flKPack, "pack", "", "exported pack")

	kernelImportFlag.Parse(args)

	if flName == "" {
		cmdEnv.Logger.Fatalln("please provide name for this new kernel")
	}

	if flVmlinux == "" && flKPack == "" {
		cmdEnv.Logger.Fatalln("Please provide vmlinux or pack path")
	}

	if flVmlinux != "" && flKPack != "" {
		cmdEnv.Logger.Fatalln("Vmlinux and Pack can not be provided together")
	}

	// Import by providing vmlinux
	if flVmlinux != "" {
		fi, err := os.Stat(flVmlinux)
		if err != nil || fi.IsDir() {
			cmdEnv.Logger.Fatalln("can not read vmlinux file.")
		}

		if flModule != "" {
			fi, err := os.Stat(flModule)
			if err != nil || !fi.IsDir() {
				cmdEnv.Logger.Fatalln("can not read module dir.")
			}
		}
	}

	// Import by providing exported kernel pack
	if flKPack != "" {
		fi, err := os.Stat(flKPack)
		if err != nil || fi.IsDir() {
			cmdEnv.Logger.Fatalln("can not read pack file.")
		}
	}

	kernelCatalog, err := cmdEnv.BaseCatalog.Dir("kernels")
	if err != nil {
		return err
	}
	if _, err := kernelCatalog.TryDir(flName); err == nil {
		cmdEnv.Logger.Fatalf("kernel %s already exists, please try another name.", flName)
	}

	// FIXME: cleanup this dir if error happens.
	newKernelCatalog, err := kernelCatalog.Dir(flName)
	if err != nil {
		cmdEnv.Logger.Fatalln("can not create new kernel catalog.")
	}

	ki := new(kernel.KernelInfo)
	ki.Name = flName
	ki.Version = flVersion

	// FIXME: support kernel pack import
	vmlinux, err := newKernelCatalog.File("vmlinux")
	if err != nil {
		cmdEnv.Logger.Fatalln("can not create vmlinux file.")
	}

	if err := util.Copy(flVmlinux, vmlinux); err != nil {
		cmdEnv.Logger.Fatalln("Copy kernel failed.")
	}

	initrdTmp, err := kernel.BuildInitrd(flModule)
	if err != nil {
		cmdEnv.Logger.Fatalf("can not create initrd: %s\n", err)
	}
	defer os.Remove(initrdTmp)

	initrd, err := newKernelCatalog.File("initrd")
	if err != nil {
		cmdEnv.Logger.Fatalln("can not create initrd file.")
	}
	if err := util.Copy(initrdTmp, initrd); err != nil {
		cmdEnv.Logger.Fatalln("Copy initrd failed.")
	}

	info, err := newKernelCatalog.File("info.json")
	if err != nil {
		cmdEnv.Logger.Fatalln("can not create info.json file.")
	}

	infof, err := os.Create(info)
	if err != nil {
		cmdEnv.Logger.Fatalln("can not create info.json file.")
	}
	defer infof.Close()

	bs, err := json.Marshal(ki)
	if err != nil {
		cmdEnv.Logger.Fatalln("can not create info.json file.")
	}
	infof.Write(bs)

	return nil
}

func cmdKernelList(args []string, cmdEnv *CommandEnv) error {
	kernelCatalog, err := cmdEnv.BaseCatalog.Dir("kernels")
	if err != nil {
		cmdEnv.Logger.WithFields(logrus.Fields{
			"module": "kernel",
		}).Fatalln("can not read kernel catalog.")
	}

	for _, catalog := range kernelCatalog.Dirs() {
		fp, err := catalog.TryFile("info.json")
		if err != nil {
			continue
		}

		b, err := ioutil.ReadFile(fp)
		if err != nil {
			cmdEnv.Logger.Debugf("Read kernel info file failed: %s\n", err)
			continue
		}

		var kInfo kernel.KernelInfo
		if err := json.Unmarshal(b, &kInfo); err != nil {
			cmdEnv.Logger.Debugf("Unmarshal failed: %s\n", err)
			continue
		}

		cmdEnv.Logger.Infof("%s %s", kInfo.Name, kInfo.Version)
	}

	return nil
}
