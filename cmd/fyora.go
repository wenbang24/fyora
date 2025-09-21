package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/gobwas/glob"
	yaml "github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
)

var Version = "v1.3.0"

type Link struct {
	Type   string `yaml:"type"`
	Source string `yaml:"source"`
	Dest   string `yaml:"target"`
	Unsafe bool   `yaml:"unsafe"`
}

type Config struct {
	Links      []Link   `yaml:"links"`
	Ignore     []string `yaml:"ignore"`
	IgnoreGlob []glob.Glob
}

var ConfigFile string
var config Config

var rootCmd = &cobra.Command{
	Use:   "fyora",
	Short: "Fyora: a declarative replacement to GNU Stow",
	Long: `Fyora is a declarative replacement to GNU Stow. It allows you to manage your dotfiles and other configuration files in a more organized and efficient way.
Made with love by @wenbang24
Docs: https://github.com/wenbang24/fyora/blob/main/README.md`,
	Version: Version,
	RunE: func(cmd *cobra.Command, args []string) error {
		ConfigFile, err := removeHomeDir(ConfigFile)
		if err != nil {
			return err
		}
		fmt.Println("Using config file:", ConfigFile)
		configFile, err := os.ReadFile(ConfigFile)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("Specified config file could not be found. You can create a default config file at by ~/.config/fyora.yaml running 'fyora init'")
				return nil
			}
			fmt.Println("Error opening config file:")
			return err
		}
		if err := yaml.Unmarshal([]byte(configFile), &config); err != nil {
			fmt.Println("Error reading config file:")
			return err
		}
		for _, ignore := range config.Ignore {
			g, err := glob.Compile(ignore)
			if err != nil {
				fmt.Printf("Error compiling glob pattern '%s': %s\n", ignore, err)
				return err
			}
			config.IgnoreGlob = append(config.IgnoreGlob, g)
		}
		count := 0
		var wg sync.WaitGroup
		for _, link := range config.Links {
			wg.Add(1)
			go func(link Link) {
				defer wg.Done()
				if link.Type == "outside" || link.Type == "file" {
					if err := OutsideSymlink(link); err != nil {
						fmt.Printf("Error creating symlink: %s\n", err)
						count--
					}
				} else if link.Type == "inside" {
					if err := InsideSymlink(link); err != nil {
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
	rootCmd.AddCommand(initCmd)
	rootCmd.Flags().StringVarP(&ConfigFile, "config", "c", "~/.config/fyora.yaml", "Path to the configuration file")
	var err error
	ConfigFile, err = removeHomeDir(ConfigFile)
	if err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			os.Exit(134) // something terribly catastrophic happened (how does printing to stderr fail tho???)
		}
		os.Exit(1)
	}
	if err := rootCmd.Execute(); err != nil {
		if _, err := fmt.Fprintln(os.Stderr, err); err != nil {
			os.Exit(134) // something terribly catastrophic happened (how does printing to stderr fail tho???)
		}
		os.Exit(1)
	}
}
