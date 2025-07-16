package cmd_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	fyora "github.com/wenbang24/fyora/cmd"
	"gopkg.in/yaml.v3"
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

func createFiles(files []string, prefix string) error {
	for _, file := range files {
		path := filepath.Join(dname, prefix, file)
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
	exit := m.Run()
	if err := os.RemoveAll(dname); err != nil {
		fmt.Println("Error removing temporary directory:", err)
		os.Exit(1)
	}
	os.Exit(exit)
}

func TestOutside(t *testing.T) {
	err := createDirs([]string{"outside", "outside/nested", "outside_target"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = createFiles([]string{"test1.txt", "nested/test2.txt"}, "outside")
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
	if string(test1) != "test1.txt qwertyuiop" {
		t.Fatalf("Test file content mismatch: expected %q, got %q\n", "test1.txt qwertyuiop", string(test1))
	}
	test2, err := os.ReadFile(filepath.Join(dname, "outside_target/outside/nested/test2.txt"))
	if err != nil {
		t.Fatalf("Failed to read nested test file: %v\n", err)
	}
	if string(test2) != "nested/test2.txt qwertyuiop" {
		t.Fatalf("Nested test file content mismatch: expected %q, got %q\n", "test2.txt qwertyuiop", string(test2))
	}
}

func TestInside(t *testing.T) {
	err := createDirs([]string{"inside", "inside/nested", "inside/nested/nested", "inside_target"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = createFiles([]string{"test1.txt", "nested/test2.txt", "nested/nested/test3.txt"}, "inside")
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
	if string(test1) != "test1.txt qwertyuiop" {
		t.Fatalf("Test file content mismatch: expected %q, got %q\n", "test1.txt qwertyuiop", string(test1))
	}
	test2, err := os.ReadFile(filepath.Join(dname, "inside_target/nested/test2.txt"))
	if err != nil {
		t.Fatalf("Failed to read nested test file: %v\n", err)
	}
	if string(test2) != "nested/test2.txt qwertyuiop" {
		t.Fatalf("Nested test file content mismatch: expected %q, got %q\n", "test2.txt qwertyuiop", string(test2))
	}
	test3, err := os.ReadFile(filepath.Join(dname, "inside_target/nested/nested/test3.txt"))
	if err != nil {
		t.Fatalf("Failed to read deeply nested test file: %v\n", err)
	}
	if string(test3) != "nested/nested/test3.txt qwertyuiop" {
		t.Fatalf("Deeply nested test file content mismatch: expected %q, got %q\n", "test3.txt qwertyuiop", string(test3))
	}
}

/*
main/
- 1.txt
- 2.txt
- .dotfile
- setup.sh
- .git
--- gitfile.txt
--- hook.sh
- dir1
--- 1.txt
--- dir2
----- dir3
------- 1.txt
------- dir3.txt
------- dir4
--------- dir4.txt
--------- 1.txt
*/

type fyoraYaml struct {
	Links  []map[string]string
	Ignore []string
}

func TestGlobIgnore(t *testing.T) {
	err := createDirs([]string{"main", "main/dir1", "main/dir1/dir2", "main/dir1/dir2/dir3", "main/dir1/dir2/dir3/dir4", ".git"})
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}
	err = createFiles([]string{"1.txt", "2.txt", ".dotfile", "setup.sh", ".git/gitfile.txt", ".git/hook.sh", "dir1/1.txt", "dir1/dir2/dir3/1.txt", "dir1/dir2/dir3/dir3.txt", "dir1/dir2/dir3/dir4/1.txt", "dir1/dir2/dir3/dir4/dir4.txt"}, "main")
	if err != nil {
		t.Fatalf("Failed to create files: %v", err)
	}
	configFilename := filepath.Join(dname, "fyora.yaml")
	configFile, err := os.OpenFile(configFilename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		t.Fatalf("Failed to open config file %s: %v", configFilename, err)
	}
	defer configFile.Close()
	err = yaml.NewEncoder(configFile).Encode(fyoraYaml{
		Links: []map[string]string{
			{
				"type":   "inside",
				"source": filepath.Join(dname, "main"),
				"dest":   filepath.Join(dname, "main_target"),
			},
		},
		Ignore: []string{
			"*.sh",
			"**/dir1/**/1.txt",
			".*",
		},
	})
	if err != nil {
		t.Fatalf("Failed to write config file %s: %v", configFilename, err)
	}
	link := fyora.Link{
		Type:   "inside",
		Source: filepath.Join(dname, "main"),
		Dest:   filepath.Join(dname, "main_target"),
	}
	if err := fyora.InsideSymlink(link); err != nil {
		t.Fatalf("Failed to create inside symlink: %v", err)
	}
}
