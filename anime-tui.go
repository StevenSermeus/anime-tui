package main

import (
	"github.com/StevenSermeus/anime-tui/cmd"
	"github.com/StevenSermeus/anime-tui/version"
)

func main() {
	cmd.Execute()
	version.CheckForUpdate()
}
