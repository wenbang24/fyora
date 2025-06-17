package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	fyora "github.com/wenbang24/fyora/cmd"
)

var dname = filepath.Join(os.TempDir(), "fyora_test")

func createDirs(dirs []string) error {
	for _, dir := range dirs {
		path := filepath.Join(dname, dir)
		if err := os.MkdirAll(path, 0777); err != nil && !os.IsExist(err) {
			return fmt.Errorf("Failed to create directory %s: %w", path, err)
		}
	}
	return nil
}

func createFiles(files []string) error {
	for _, file := range files {
		path := filepath.Join(dname, file)
		if err := os.WriteFile(path, []byte(file+" qwertyuiop"), 0644); err != nil && !os.IsExist(err) {
			return fmt.Errorf("Failed to create file %s: %w", path, err)
		}
	}
	return nil
}

func TestMain(m *testing.M) {
	err := os.RemoveAll(dname)

	err = os.Mkdir(dname, 0777)
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
	err := createDirs([]string{"outside", "outside/nested", "outside_target"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = createFiles([]string{"outside/test1.txt", "outside/nested/test2.txt"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	link := fyora.Link{
		Type:   "outside",
		Source: filepath.Join(dname, "outside"),
		Dest:   filepath.Join(dname, "outside_target"),
	}
	if err := fyora.OutsideSymlink(link); err != nil {
		t.Fatalf("Failed to create outside symlink: %v", err)
	}
	test1, err := os.ReadFile(filepath.Join(dname, "outside_target/outside/test1.txt"))
	if err != nil {
		t.Fatalf("Failed to read test file: %v\n", err)
	}
	if string(test1) != "outside/test1.txt qwertyuiop" {
		t.Fatalf("Test file content mismatch: expected %q, got %q\n", "test1.txt qwertyuiop", string(test1))
	}
	test2, err := os.ReadFile(filepath.Join(dname, "outside_target/outside/nested/test2.txt"))
	if err != nil {
		t.Fatalf("Failed to read nested test file: %v\n", err)
	}
	if string(test2) != "outside/nested/test2.txt qwertyuiop" {
		t.Fatalf("Nested test file content mismatch: expected %q, got %q\n", "test2.txt qwertyuiop", string(test2))
	}
}

func TestInside(t *testing.T) {
	err := createDirs([]string{"inside", "inside/nested", "inside/nested/nested", "inside_target"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = createFiles([]string{"inside/test1.txt", "inside/nested/test2.txt", "inside/nested/nested/test3.txt"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	link := fyora.Link{
		Type:   "inside",
		Source: filepath.Join(dname, "inside"),
		Dest:   filepath.Join(dname, "inside_target"),
	}
	if err := fyora.InsideSymlink(link); err != nil {
		t.Fatalf("Failed to create inside symlink: %v", err)
	}
	test1, err := os.ReadFile(filepath.Join(dname, "inside_target/test1.txt"))
	if err != nil {
		t.Fatalf("Failed to read test file: %v\n", err)
	}
	if string(test1) != "inside/test1.txt qwertyuiop" {
		t.Fatalf("Test file content mismatch: expected %q, got %q\n", "test1.txt qwertyuiop", string(test1))
	}
	test2, err := os.ReadFile(filepath.Join(dname, "inside_target/nested/test2.txt"))
	if err != nil {
		t.Fatalf("Failed to read nested test file: %v\n", err)
	}
	if string(test2) != "inside/nested/test2.txt qwertyuiop" {
		t.Fatalf("Nested test file content mismatch: expected %q, got %q\n", "test2.txt qwertyuiop", string(test2))
	}
	test3, err := os.ReadFile(filepath.Join(dname, "inside_target/nested/nested/test3.txt"))
	if err != nil {
		t.Fatalf("Failed to read deeply nested test file: %v\n", err)
	}
	if string(test3) != "inside/nested/nested/test3.txt qwertyuiop" {
		t.Fatalf("Deeply nested test file content mismatch: expected %q, got %q\n", "test3.txt qwertyuiop", string(test3))
	}
}
