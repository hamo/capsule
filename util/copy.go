package util

import (
	"io"
	"os"
)

func Copy(src string, dst string) error {
	// FIXME: support DIR
	dstf, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstf.Close()

	srcf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcf.Close()

	_, err = io.Copy(dstf, srcf)
	if err != nil {
		return err
	}
	return nil
}
