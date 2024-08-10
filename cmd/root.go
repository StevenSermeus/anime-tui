package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anime-tui",
	Short: "A terminal user interface for anime",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("This is the root command")
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of anime-tui",
	Run: func(cmd *cobra.Command, args []string) {
		// Only change when the tag is changed
		// When in othen branch than main put as -dev
		fmt.Println("0.0.1-beta")
	},
}

func Execute() {
	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
