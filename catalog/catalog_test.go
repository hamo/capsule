package catalog

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestEmptyBaseCatalogDir(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Create temp dir failed: %s\n", err)
	}
	defer os.RemoveAll(tmpdir)

	cd, err := NewBaseCatalogDir(tmpdir, "baseDir")
	if err != nil {
		t.Fatalf("Create base catalog dir failed: %s", err)
	}

	if cd.Name != "baseDir" ||
		cd.Parent != nil ||
		len(cd.subDir) != 0 ||
		len(cd.subFile) != 0 {
		t.Fatalf("Empty catalog dir created with contents: %+v", cd)
	}

	childDir, err := cd.Dir("child1")
	if err != nil {
		t.Fatalf("create child dir failed: %s", err)
	}

	fi, err := os.Stat(childDir.Path)
	if err != nil || !fi.IsDir() {
		t.Fatalf("Create child dir failed.")
	}
}

func TestExistBaseCatalogDir(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("Create temp dir failed: %s\n", err)
	}
	defer os.RemoveAll(tmpdir)

	baseDir := filepath.Join(tmpdir, "baseDir")
	childDir1 := filepath.Join(baseDir, "child1")
	if err := os.MkdirAll(childDir1, 0755); err != nil {
		t.Fatalf("Create child dir by os.Mkdir failed: %s", err)
	}

	childDir2 := filepath.Join(baseDir, "child2")
	if err := os.MkdirAll(childDir2, 0755); err != nil {
		t.Fatalf("Create child dir by os.Mkdir failed: %s", err)
	}

	cd, err := NewBaseCatalogDir(tmpdir, "baseDir")
	if err != nil {
		t.Fatalf("Create base catalog dir failed: %s", err)
	}

	if cd.Name != "baseDir" ||
		cd.Parent != nil ||
		len(cd.subDir) != 2 ||
		len(cd.subFile) != 0 {
		t.Fatalf("Exist catalog dir created with contents: %+v", cd)
	}

	childDir, err := cd.Dir("child1")
	if err != nil {
		t.Fatalf("create child dir failed: %s", err)
	}

	fi, err := os.Stat(childDir.Path)
	if err != nil || !fi.IsDir() {
		t.Fatalf("Create child dir failed.")
	}
}
