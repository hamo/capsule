package kernel

import (
	"io"
	"io/ioutil"
	// FIXME: dynamic linking
	"bytes"
	"compress/gzip"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	busyBoxBaseURL = "http://busybox.net/downloads/binaries"
	busyBoxVersion = "1.21.1"
	busyBoxName    = "busybox-x86_64"

	busyBoxURL = busyBoxBaseURL + "/" + busyBoxVersion + "/" + busyBoxName
)

func BuildInitrd(moduleDir string) (string, error) {
	initrdDir, err := ioutil.TempDir("", "initrd")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(initrdDir)

	// Create dirs under initrdDir
	dirs := [][]string{
		[]string{initrdDir, "sysroot"},
		[]string{initrdDir, "bin"},
		[]string{initrdDir, "sbin"},
		[]string{initrdDir, "proc"},
		[]string{initrdDir, "sys"},
		[]string{initrdDir, "run"},
		[]string{initrdDir, "dev"},
		[]string{initrdDir, "etc", "modprobe.d"},
		[]string{initrdDir, "lib", "modules"},
	}
	for _, dir := range dirs {
		// FIXME: filepath.Join(initrdDir, dir...)
		if err := os.MkdirAll(filepath.Join(dir...), 0755); err != nil {
			return "", err
		}
	}

	// Download busybox and create symlink
	busyboxPath := filepath.Join([]string{initrdDir, "bin", "busybox"}...)
	if err := downloadURL(busyBoxURL, busyboxPath); err != nil {
		return "", err
	}
	if err := os.Chmod(busyboxPath, 0755); err != nil {
		return "", err
	}

	binFiles := []string{
		"cat",
		"cp",
		"dmesg",
		"echo",
		"ls",
		"lsmod",
		"mkdir",
		"rm",
		"sh",
		"sleep",
	}

	for _, v := range binFiles {
		if err := os.Symlink("./busybox", filepath.Join(initrdDir, "bin", v)); err != nil {
			return "", err
		}
	}

	sbinFiles := []string{
		"insmod",
		"rmmod",
		"mdev",
		"mknod",
		"modprobe",
		"mount",
		"switch_root",
		"umount",
	}

	for _, v := range sbinFiles {
		if err := os.Symlink("../bin/busybox", filepath.Join(initrdDir, "sbin", v)); err != nil {
			return "", err
		}
	}

	// copy all modules into initrd
	// FIXME: review, this will make initrd too big to be loaded by a small memory system
	// FIXME: move to a golang native solution
	// FIXME: hardlink?
	if moduleDir != "" {
		cp := exec.Command("cp", "-rL", moduleDir, filepath.Join(initrdDir, "lib", "modules"))
		if err := cp.Run(); err != nil {
			return "", err
		}
	}

	// copy /init
	initFile, err := os.Create(filepath.Join(initrdDir, "init"))
	if err != nil {
		return "", err
	}
	if _, err := io.WriteString(initFile, initShell); err != nil {
		return "", err
	}
	if err := initFile.Chmod(0755); err != nil {
		return "", err
	}
	initFile.Close()

	// create cpio
	// FIXME: move to a golang native solution
	var fileList bytes.Buffer
	findCmd := exec.Command("find", ".")
	findCmd.Dir = initrdDir
	findCmd.Stdout = &fileList
	if err := findCmd.Run(); err != nil {
		return "", err
	}

	var initrdBuffer bytes.Buffer

	cpioCmd := exec.Command("cpio", "-H", "newc", "-o")
	cpioCmd.Dir = initrdDir
	cpioCmd.Stdin = &fileList
	cpioCmd.Stdout = &initrdBuffer
	if err := cpioCmd.Run(); err != nil {
		return "", err
	}

	gzipInitrdFile, err := ioutil.TempFile("", "initrd")
	if err != nil {
		return "", err
	}
	gzipInitrdWriter := gzip.NewWriter(gzipInitrdFile)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(gzipInitrdWriter, &initrdBuffer)
	if err != nil {
		return "", err
	}
	if err := gzipInitrdWriter.Flush(); err != nil {
		return "", err
	}
	gzipInitrdWriter.Close()
	gzipInitrdFile.Close()

	return gzipInitrdFile.Name(), nil
}

func downloadURL(url string, output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()

	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return err
	}
	return nil
}
