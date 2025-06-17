package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	Errors "github.com/pkg/errors"
)

func OutsideSymlink(link Link) error {
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
	filename := filepath.Base(source)
	dest = filepath.Join(dest, filename)
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

func InsideSymlink(link Link) error {
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
		if _, skip := config.IgnoreSet[file.Name()]; skip {
			continue
		}
		link := Link{
			Type:   "outside",
			Source: filepath.Join(sourceDir, file.Name()),
			Dest:   filepath.Join(destDir, file.Name()),
			Unsafe: link.Unsafe,
		}
		if err := OutsideSymlink(link); err != nil {
			fmt.Printf("Error creating symlink for %s to %s: %s\n", filepath.Join(sourceDir, file.Name()), filepath.Join(destDir, file.Name()), err)
			continue
		}
	}
	return nil
}
