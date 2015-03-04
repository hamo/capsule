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

func TestCleanup(t *testing.T) {
	/*
	*  base -> dir11 -> dir21          -> dir31         -> file41(clean)
	*                                 |-> file31(keep)
	*               |-> file21(clean)
	*               |-> file22(keep)
	*               |-> dir22          -> dir32         -> file42(keep)
	*                                 |-> file32(clean)
	 */

	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("create working dir failed: %s\n", err)
	}
	defer os.RemoveAll(tmpdir)

	dir11, err := NewBaseCatalogDir(tmpdir, "dir11")
	if err != nil {
		t.Fatalf("create base catalog failed: %s\n", err)
	}

	dir21, err := dir11.Dir("dir21")
	if err != nil {
		t.Fatalf("create dir21 catalog failed: %s\n", err)
	}
	dir22, err := dir11.Dir("dir22")
	if err != nil {
		t.Fatalf("create dir22 catalog failed: %s\n", err)
	}
	file21, err := dir11.File("file21", true)
	if err != nil {
		t.Fatalf("create file21 catalog failed: %s\n", err)
	}
	f21, err := os.Create(file21)
	if err != nil {
		t.Fatalf("create file21 catalog failed: %s\n", err)
	}
	f21.Close()
	file22, err := dir11.File("file22", false)
	if err != nil {
		t.Fatalf("create file22 catalog failed: %s\n", err)
	}
	f22, err := os.Create(file22)
	if err != nil {
		t.Fatalf("create file22 catalog failed: %s\n", err)
	}
	f22.Close()

	file31, err := dir21.File("file31", false)
	if err != nil {
		t.Fatalf("create file31 catalog failed: %s\n", err)
	}
	f31, err := os.Create(file31)
	if err != nil {
		t.Fatalf("create file31 catalog failed: %s\n", err)
	}
	f31.Close()
	dir31, err := dir21.Dir("dir31")
	if err != nil {
		t.Fatalf("create dir31 catalog failed: %s\n", err)
	}

	file32, err := dir22.File("file32", true)
	if err != nil {
		t.Fatalf("create file32 catalog failed: %s\n", err)
	}
	f32, err := os.Create(file32)
	if err != nil {
		t.Fatalf("create file32 catalog failed: %s\n", err)
	}
	f32.Close()
	dir32, err := dir22.Dir("dir32")
	if err != nil {
		t.Fatalf("create dir32 catalog failed: %s\n", err)
	}

	file41, err := dir31.File("file41", true)
	if err != nil {
		t.Fatalf("create file41 catalog failed: %s\n", err)
	}
	f41, err := os.Create(file41)
	if err != nil {
		t.Fatalf("create file41 catalog failed: %s\n", err)
	}
	f41.Close()

	file42, err := dir32.File("file42", false)
	if err != nil {
		t.Fatalf("create file42 catalog failed: %s\n", err)
	}
	f42, err := os.Create(file42)
	if err != nil {
		t.Fatalf("create file42 catalog failed: %s\n", err)
	}
	f42.Close()

	// create test file structure done
	dir11.Cleanup(false)

	// file42 should exist
	if _, err := os.Stat(file42); err != nil {
		t.Fatalf("can not access file42 after cleanup: %s\n", err)
	}

	// dir31 should disappear
	if _, err := os.Stat(dir31.Path); err == nil {
		t.Fatalf("still can access dir31 after cleanup")
	}

	// file32 should disappear
	if _, err := os.Stat(file32); err == nil {
		t.Fatalf("still can access file32 after cleanup")
	}

	// file31 should exist
	if _, err := os.Stat(file31); err != nil {
		t.Fatalf("can not access file31 after cleanup: %s\n", err)
	}

	// dir21 and dir22 should exist
	if _, err := os.Stat(dir21.Path); err != nil {
		t.Fatalf("can not access dir21 after cleanup: %s\n", err)
	}
	if _, err := os.Stat(dir22.Path); err != nil {
		t.Fatalf("can not access dir22 after cleanup: %s\n", err)
	}

	// file21 should disapper
	if _, err := os.Stat(file21); err == nil {
		t.Fatalf("still can access file21 after cleanup")
	}
	// file22 should exist
	if _, err := os.Stat(file22); err != nil {
		t.Fatalf("can not access file22 after cleanup: %s\n", err)
	}

	// then force cleanup
	dir11.Cleanup(true)

	// dir11 should disappear
	if _, err := os.Stat(dir11.Path); err == nil {
		t.Fatalf("still can access dir11 after force cleanup.")
	}
}
