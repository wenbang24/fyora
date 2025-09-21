package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var defaultConfig = `# Fyora config file
# Docs: https://github.com/wenbang24/fyora/blob/main/README.md
links:
    # Add your links here
    # Example:
    # - type: outside
    #   source: /dir1
    #   target: ~/dir2
    #   unsafe: true
    # - type: inside
    #   source: ~/dir3
    #   target: ~/dir/dir4
    # - type: file
    #   source: /dir5/file.txt
    #   target: ~/dir2/dir/file.txt
ignore:
    - .DS_Store
    - .git
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a default config file at ~/.config/fyora.yaml",
	Long: `Initialize a default config file at ~/.config/fyora.yaml.
If a config file already exists, it will not be overwritten.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := removeHomeDir("~/.config/fyora.yaml")
		if err != nil {
			return err
		}
		_, err = os.Stat(configPath)
		if os.IsNotExist(err) {
			err = os.WriteFile(configPath, []byte(defaultConfig), 0644)
			if err != nil {
				fmt.Println("Error creating config file at ~/.config/fyora.yaml:")
				return err
			}
			fmt.Println("Default config file created at ~/.config/fyora.yaml")
			return nil
		} else if err == nil {
			fmt.Println("Config file already exists at ~/.config/fyora.yaml. Not overwriting.")
			return nil
		}
		fmt.Println("Error checking config file at ~/.config/fyora.yaml:")
		return err
	},
}
