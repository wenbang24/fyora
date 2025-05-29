package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	yaml "github.com/goccy/go-yaml"
	Errors "github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type Link struct {
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	Dest   string `yaml:"target"`
	Unsafe bool   `yaml:"unsafe"`
}

type Config struct {
	Links     []Link   `yaml:"links"`
	Ignore    []string `yaml:"ignore"`
	IgnoreSet map[string]struct{}
}

var ConfigFile string

func isSymlink(path string) (bool, error) {
	file, err := os.Lstat(path)
	if Errors.Is(err, os.ErrNotExist) {
		fmt.Printf("Creating symlink from %s\n", path)
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return file.Mode()&os.ModeSymlink != 0, nil
}

func removeHomeDir(path string) (string, error) {
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", homeDir, 1)
	}
	return path, nil
}

func pathType(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return "dne", nil
		}
		return "", err
	}
	if info.IsDir() {
		return "directory", nil
	}
	return "file", nil
}

func outsideSymlink(link Link) error {
	source, err := removeHomeDir(link.Source)
	if err != nil {
		fmt.Println("Error getting absolute path of source:")
		return err
	}
	dest, err := removeHomeDir(link.Dest)
	if err != nil {
		fmt.Println("Error getting absolute path of target:")
		return err
	}
	sourceType, err := pathType(source)
	if err != nil {
		fmt.Println("Error checking source type:")
		return err
	}
	if sourceType == "file" {
		filename := filepath.Base(source)
		if !strings.HasSuffix(dest, filename) {
			dest = filepath.Join(dest, filename)
		}
	}
	symlink, err := isSymlink(dest)
	if err != nil {
		fmt.Println("Error checking if target is a symlink:")
		return err
	}
	if symlink {
		target, err := filepath.EvalSymlinks(dest)
		if err != nil {
			fmt.Println("Error evaluating symlink:")
			return err
		}
		if target == source {
			fmt.Printf("Symlink %s already exists and points to %s\n", dest, source)
			return nil
		} else {
			fmt.Printf("Symlink %s already exists and points to %s", dest, target)
			return Errors.New("symlink already exists and points to a different target")
		}
	}
	if link.Unsafe {
		info, err := os.Stat(dest)
		if err != nil && !os.IsNotExist(err) {
			fmt.Printf("Failed to check target: %s\n", err)
			return Errors.New("failed to check target")
		}
		if info.IsDir() {
			if err := os.RemoveAll(dest); err != nil {
				fmt.Printf("Failed to delete directory %q: %v\n", dest, err)
			} else {
				fmt.Println("Deleted directory:", dest)
			}
		} else {
			if err := os.Remove(dest); err != nil {
				fmt.Printf("Failed to delete file %q: %v\n", dest, err)
			} else {
				fmt.Println("Deleted file:", dest)
			}
		}
	}
	if err := os.Symlink(source, dest); err != nil {
		fmt.Println("Error creating symlink:")
		return err
	}
	return nil
}

func insideSymlink(link Link, ignoreSet map[string]struct{}) error {
	sourceDir, err := removeHomeDir(link.Source)
	if err != nil {
		fmt.Println("Error getting absolute path of source:")
		return err
	}
	source, err := os.Stat(sourceDir)
	if err != nil {
		fmt.Println("Error reading target directory:")
		return err
	}
	if !source.IsDir() {
		fmt.Printf("Source %s is not a directory\n", source)
		return Errors.New("source is not a directory")
	}
	destDir, err := removeHomeDir(link.Dest)
	if err != nil {
		fmt.Println("Error getting absolute path of target:")
		return err
	}
	dest, err := os.Stat(destDir)
	if err != nil {
		fmt.Println("Error reading target directory:")
		return err
	}
	if !dest.IsDir() {
		fmt.Printf("Target %s is not a directory\n", dest)
		return Errors.New("target is not a directory")
	}
	files, err := os.ReadDir(sourceDir)
	if err != nil {
		fmt.Println("Error reading source directory:")
		return err
	}
	for _, file := range files {
		if _, skip := ignoreSet[file.Name()]; skip {
			continue
		}
		link := Link{
			Type:   "outside",
			Source: filepath.Join(sourceDir, file.Name()),
			Dest:   filepath.Join(destDir, file.Name()),
			Unsafe: link.Unsafe,
		}
		if err := outsideSymlink(link); err != nil {
			fmt.Printf("Error creating symlink for %s to %s: %s\n", filepath.Join(sourceDir, file.Name()), filepath.Join(destDir, file.Name()), err)
			continue
		}
	}
	return nil
}

var rootCmd = &cobra.Command{
	Use:   "fyora",
	Short: "Fyora: a declarative replacement to GNU Stow",
	Long: `Fyora is a declarative replacement to GNU Stow.
It allows you to manage your dotfiles and other configuration files in a more organized and efficient way.
Made with love by @wenbang24`,
	RunE: func(cmd *cobra.Command, args []string) error {
		config := Config{}
		configFile, err := os.ReadFile(ConfigFile)
		if err != nil {
			fmt.Println("Error opening config file:")
			return err
		}
		if err := yaml.Unmarshal([]byte(configFile), &config); err != nil {
			fmt.Println("Error reading config file:")
			return err
		}
		var ignoreSet = make(map[string]struct{})
		for _, ignore := range config.Ignore {
			ignoreSet[ignore] = struct{}{}
		}
		count := 0
		var wg sync.WaitGroup
		for _, link := range config.Links {
			wg.Add(1)
			go func(link Link) {
				defer wg.Done()
				if link.Type == "outside" || link.Type == "file" {
					if err := outsideSymlink(link); err != nil {
						fmt.Printf("Error creating symlink: %s\n", err)
						count--
					}
				} else if link.Type == "inside" {
					if err := insideSymlink(link, ignoreSet); err != nil {
						fmt.Printf("Error creating symlink: %s\n", err)
						count--
					}
				} else {
					fmt.Printf("Unknown type: %s\n", link.Type)
					count--
				}
				count++
			}(link)
		}
		wg.Wait()
		fmt.Printf("Created %d symlinks\n", count)
		return nil
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&ConfigFile, "config", "c", "~/.config/fyora.yaml", "Path to the configuration file")
	if err := rootCmd.Execute(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			os.Exit(134) // something terribly catastrophic happened (how does printing to stderr fail tho???)
		}
		os.Exit(1)
	}
}
