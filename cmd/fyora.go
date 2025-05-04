package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "fyora",
	Short: "Fyora: a declarative replacement to GNU Stow",
	Long: `Fyora is a declarative replacement to GNU Stow. It allows you to manage your dotfiles and other configuration files in a more organized and efficient way.
Made with love by @wenbang24`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
