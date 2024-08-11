package cmd

import (
	"fmt"
	"os"

	"github.com/StevenSermeus/anime-tui/tui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "anime-tui",
	Short: "A terminal user interface for anime",
	Run: func(cmd *cobra.Command, args []string) {
		tui.StartTui()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of anime-tui",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	rootCmd.AddCommand(versionCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
