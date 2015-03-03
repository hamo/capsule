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

	subDir map[string]*CatalogDir

	// whether cleanup file
	subFile map[string]bool
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
		cd.subFile = make(map[string]bool)
		return cd.Sync()
	} else if err == nil && !fi.IsDir() {
		return nil, fmt.Errorf("%s exists but not a dir.", fp)
	} else {
		if err := os.Mkdir(fp, 0755); err != nil {
			return nil, err
		}
		return &CatalogDir{
			Name:    name,
			Path:    fp,
			Parent:  nil,
			subDir:  make(map[string]*CatalogDir),
			subFile: make(map[string]bool),
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
			d.subFile[fi.Name()] = false
		}
	}

	return d, nil
}

func (d *CatalogDir) TryDir(name string) (*CatalogDir, error) {
	cd, ok := d.subDir[name]
	if ok {
		return cd, nil
	}

	fp := filepath.Join(d.Path, name)
	fi, err := os.Stat(fp)

	if err != nil || !fi.IsDir() {
		return nil, errors.New("Dir does not exists.")
	}

	cd = &CatalogDir{
		Name:    name,
		Parent:  d,
		Path:    fp,
		subDir:  make(map[string]*CatalogDir),
		subFile: make(map[string]bool),
	}

	if _, err := cd.Sync(); err != nil {
		return nil, err
	}
	d.subDir[name] = cd
	return cd, nil
}

func (d *CatalogDir) Dir(name string) (*CatalogDir, error) {
	cd, err := d.TryDir(name)
	if err == nil {
		return cd, err
	}

	fp := filepath.Join(d.Path, name)

	if err := os.Mkdir(fp, 0755); err != nil {
		return nil, err
	}

	if _, ok := d.subDir[name]; ok {
		panic("")
	}

	cd = &CatalogDir{
		Name:    name,
		Path:    fp,
		Parent:  d,
		subDir:  make(map[string]*CatalogDir),
		subFile: make(map[string]bool),
	}

	d.subDir[name] = cd
	return cd, nil
}

func (d *CatalogDir) Dirs() map[string]*CatalogDir {
	return d.subDir
}

func (d *CatalogDir) TryFile(name string) (string, error) {
	fp := filepath.Join(d.Path, name)

	if _, ok := d.subFile[name]; ok {
		return fp, nil
	}

	fi, err := os.Stat(fp)

	if err != nil || fi.IsDir() {
		return "", errors.New("File does not exist.")
	}

	// FIXME: cleanup
	d.subFile[name] = false
	return fp, nil
}

func (d *CatalogDir) File(name string) (string, error) {
	if fp, err := d.TryFile(name); err == nil {
		return fp, nil
	}

	fp := filepath.Join(d.Path, name)

	// FIXME: cleanup
	d.subFile[name] = false
	return fp, nil
}

// FIXME: Clean
