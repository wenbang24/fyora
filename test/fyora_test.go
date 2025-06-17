package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	fyora "github.com/wenbang24/fyora/cmd"
)

var dname = filepath.Join(os.TempDir(), "fyora_test")

func TestMain(m *testing.M) {
	err := os.Mkdir(dname, 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Error creating temporary directory:", err)
		os.Exit(1)
	}
	defer os.RemoveAll(dname)
	in, err := os.ReadFile("fyora.yaml")
	if err != nil {
		fmt.Println("Error reading test configuration file:", err)
		os.Exit(1)
	}
	if err := os.WriteFile(filepath.Join(dname, "fyora.yaml"), in, 0644); err != nil {
		fmt.Println("Error writing test configuration file:", err)
		os.Exit(1)
	}

	fyora.ConfigFile = filepath.Join(os.TempDir(), "fyora.yaml")
	exit := m.Run()
	os.Exit(exit)
}

func TestOutside(t *testing.T) {
	err := os.Mkdir(filepath.Join(dname, "outside"), 0777)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	err = os.WriteFile(filepath.Join(dname, "outside/test.txt"), []byte("abcdef"), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	err = os.Mkdir(filepath.Join(dname, "outside_target"), 0777)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Failed to create target directory: %v", err)
	}
	link := fyora.Link{
		Type:   "outside",
		Source: filepath.Join(dname, "outside"),
		Dest:   filepath.Join(dname, "outside_target"),
	}
	if err := fyora.OutsideSymlink(link); err != nil {
		t.Fatalf("Failed to create outside symlink: %v", err)
	}
	if _, err := os.Lstat(filepath.Join(dname, "outside_target/outside/text.txt")); err != nil {
		t.Fatalf("Failed to verify symlink: %v", err)
	}
}
