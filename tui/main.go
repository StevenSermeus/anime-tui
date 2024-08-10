package tui

import (
	"fmt"
	"os"

	stplatform "github.com/StevenSermeus/anime-tui/st_platform"
)

func StartTui() {
	mav := &stplatform.Mavanimes{}
	state := "menu"
	quit := false
	anime_state := stplatform.Anime{}
	for !quit {
		switch state {
		case "menu":
			state, quit = Menu()
		case "Recherche":
			anime, isRetunring, isQuiting, err := Search(mav)
			if err != nil {
				fmt.Println("Error searching anime:", err)
				os.Exit(1)
			}
			if isQuiting {
				quit = true
				continue
			}
			if isRetunring {
				state = "menu"
				continue
			}
			anime_state = anime
			state = "Episodes"
		case "Nouveautés":
			isRetunring, isQuiting, err := Newest(mav)
			if err != nil {
				fmt.Println("Error getting newest anime:", err)
				os.Exit(1)
			}
			if isRetunring {
				state = "menu"
				continue
			}
			quit = isQuiting
		case "Episodes":
			isRetunring, isQuiting, err := Episodes(mav, anime_state)
			if err != nil {
				fmt.Println("Error getting episodes:", err)
				os.Exit(1)
			}
			if isRetunring {
				state = "Recherche"
				continue
			}
			quit = isQuiting
		}

	}
	os.Exit(0)
}
