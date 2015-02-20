package catalog

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type CatalogDir struct {
	Name string
	Path string

	Parent *CatalogDir

	subDir  map[string]*CatalogDir
	subFile []string
}

func NewBaseCatalogDir(baseDir string, name string) (*CatalogDir, error) {
	fi, err := os.Stat(baseDir)
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, errors.New("baseDir is not dir.")
	}

	fp := filepath.Join(baseDir, name)

	cd := new(CatalogDir)
	// If fp already exists, read its contents.
	fi, err = os.Stat(fp)
	if err == nil && fi.IsDir() {
		cd.Name = name
		cd.Parent = nil
		cd.Path = fp
		cd.subDir = make(map[string]*CatalogDir)
		return cd.Sync()
	} else if err == nil && !fi.IsDir() {
		return nil, fmt.Errorf("%s exists but not a dir.", fp)
	} else {
		if err := os.Mkdir(fp, 0755); err != nil {
			return nil, err
		}
		return &CatalogDir{
			Name:   name,
			Path:   fp,
			Parent: nil,
			subDir: make(map[string]*CatalogDir),
		}, nil
	}
}

func (d *CatalogDir) Sync() (*CatalogDir, error) {
	fi, err := os.Stat(d.Path)
	if err != nil || !fi.IsDir() {
		return nil, fmt.Errorf("called Sync with path %s.", d.Path)
	}

	fis, err := ioutil.ReadDir(d.Path)
	if err != nil {
		return nil, err
	}

	for _, fi := range fis {
		if fi.IsDir() {
			d.Dir(fi.Name())
		} else {
			d.subFile = append(d.subFile, fi.Name())
		}
	}

	return d, nil
}

func (d *CatalogDir) Dir(name string) (*CatalogDir, error) {
	fp := filepath.Join(d.Path, name)
	fi, err := os.Stat(fp)

	cd := new(CatalogDir)

	if err == nil && fi.IsDir() {
		cd.Name = name
		cd.Parent = d
		cd.Path = fp
		cd.subDir = make(map[string]*CatalogDir)
		if _, err := cd.Sync(); err != nil {
			return nil, err
		}
		d.subDir[name] = cd
		return cd, nil
	} // FIXME: err == nil && !fi.IsDir()
	if err := os.Mkdir(fp, 0755); err != nil {
		return nil, err
	}

	cd.Name = name
	cd.Path = fp
	cd.Parent = d
	cd.subDir = make(map[string]*CatalogDir)

	if _, ok := d.subDir[name]; ok {
		panic("")
	}
	d.subDir[name] = cd
	return cd, nil
}

func (d *CatalogDir) File(name string) (string, error) {
	fp := filepath.Join(d.Path, name)
	fi, err := os.Stat(fp)

	if err == nil && !fi.IsDir() {
		d.subFile = append(d.subFile, name)
		return fp, nil
	} // FIXME: err == nil && fi.IsDir()

	d.subFile = append(d.subFile, name)
	return fp, nil
}

// FIXME: Clean
