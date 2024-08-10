package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/StevenSermeus/anime-tui/network"
	"github.com/StevenSermeus/anime-tui/tui"
	"github.com/StevenSermeus/anime-tui/video"
	"github.com/briandowns/spinner"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//This is bad code, but it's a proof of concept
	state := "menu"
	choice := tui.Item{}
	episode_choice := tui.Item{}
	decoded_xsrf := ""
	mav_token := ""
	for {
		switch state {
		case "menu":
			choice := optionList()
			if choice.Name == "Nouveaux épisodes" {
				state = "newEpisodes"
			}
			if choice.Name == "Recherche" {
				state = "search"
			}
		case "newEpisodes":
			choice, decoded_xsrf, mav_token = newEpisodes()
			if choice.Name == "Retour" {
				state = "menu"
				continue
			}
			play(choice.Url, decoded_xsrf, mav_token)
		case "search":
			choice, decoded_xsrf, mav_token = search()
			if choice.Name == "Retour" {
				state = "menu"
			} else {
				state = "episodeList"
			}
		case "episodeList":
			episode_choice = episodeList(choice)
			if episode_choice.Name == "Retour" {
				state = "search"
				continue
			}
			play(episode_choice.Url, decoded_xsrf, mav_token)
		}
	}
}

func play(url string, decoded_xsrf string, mav_token string) {
	result, err := network.GetVideoLink(url, mav_token, decoded_xsrf)
	if err != nil {
		panic(err)
	}
	link := ""
	d0000d := regexp.MustCompile(`https://d0000d.com`)
	for _, res := range result {
		if d0000d.MatchString(res) {
			link = res
		}
	}
	if link == "" {
		panic("Couldn't find the video link")
	}
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	new_Link, _ := getD0000DLink1(link)
	play_link, err := getD0000DLink2("https://d0000d.com" + new_Link)
	s.Stop()
	if err != nil {
		panic(err)
	}
	player := video.VLC{}
	player.Play(play_link, "Anime TUI")
}

func getD0000DLink2(url string) (string, error) {
	//Sleep 1 second to avoid getting blocked
	time.Sleep(1 * time.Second)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-GB,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Cookie", "lang=1;")
	req.Header.Set("DNT", "1")
	req.Header.Set("Priority", "u=0, i")
	req.Header.Set("Sec-CH-UA", `"Chromium";v="127", "Not)A;Brand";v="99"`)
	req.Header.Set("Sec-CH-UA-Mobile", "?0")
	req.Header.Set("Sec-CH-UA-Platform", `"macOS"`)
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "none")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	match := regexp.MustCompile(`href="(.*?)"`).FindAllString(string(body), -1)
	video_link := ""
	for _, url := range match {
		if regexp.MustCompile(`mp4`).MatchString(url) {
			video_link = url
		}
	}
	if video_link == "" {
		panic("Couldn't find the video link")
	}
	video_link = regexp.MustCompile(`href="`).ReplaceAllString(video_link, "")
	video_link = video_link[:len(video_link)-1]
	return video_link, nil
}

func getD0000DLink1(url string) (string, error) {
	url = regexp.MustCompile(`/e/`).ReplaceAllString(url, "/d/")
	res, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	additional_link := ""
	doc.Find("div.download-content a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		additional_link = link
	})
	return additional_link, nil
}

func episodeList(choice tui.Item) tui.Item {
	episodeList, err := network.GetEpisodeList(choice.Url)
	if err != nil {
		panic(err)
	}
	episode_list := []list.Item{}
	for _, episode := range episodeList {
		episode_list = append(episode_list, tui.Item{Name: episode[2], Url: episode[1]})
	}
	episode_list = append(episode_list, tui.Item{Name: "Retour", Url: "Retour à la recherche"})
	m := tui.Model{List: list.New(episode_list, list.NewDefaultDelegate(), 0, 0)}
	m.List.Title = "Liste des épisodes"

	p := tea.NewProgram(m, tea.WithAltScreen())

	mf, err := p.Run()
	if err != nil {
		panic(err)
	}
	m = mf.(tui.Model)
	if m.Quiting {
		os.Exit(0)
	}
	return m.Choice
}

func optionList() (choice tui.Item) {
	options := []list.Item{
		tui.Item{Name: "Nouveaux épisodes", Url: "Voir les derniers épisodes sortis"},
		tui.Item{Name: "Recherche", Url: "Rechercher un anime par nom"},
		tui.Item{Name: "Exit", Url: "Press Ctrl+C to quit"},
	}

	m := tui.Model{List: list.New(options, list.NewDefaultDelegate(), 0, 0)}
	m.List.Title = "Anime TUI"

	p := tea.NewProgram(m, tea.WithAltScreen())

	mf, err := p.Run()
	if err != nil {
		panic(err)
	}

	m = mf.(tui.Model)
	if m.Quiting || m.Choice.Name == "Exit" {
		os.Exit(0)
		return
	}
	return m.Choice
}

func newEpisodes() (choice tui.Item, decoded_xsrf string, mav_token string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	animeList, decoded_xsrf, mav_token, err := network.GetRecentEpisode()
	s.Stop()
	if err != nil {
		panic(err)
	}
	episode_list := []list.Item{}
	for _, anime := range animeList {
		episode_list = append(episode_list, tui.Item{Name: anime[1], Url: anime[0]})
	}
	episode_list = append(episode_list, tui.Item{Name: "Retour", Url: "Retour au menu"})
	m := tui.Model{List: list.New(episode_list, list.NewDefaultDelegate(), 0, 0)}
	m.List.Title = "Derniers épisodes"

	p := tea.NewProgram(m, tea.WithAltScreen())

	mf, err := p.Run()
	if err != nil {
		panic(err)
	}
	m = mf.(tui.Model)
	if m.Quiting {
		os.Exit(0)
	}
	return m.Choice, decoded_xsrf, mav_token
}

func search() (choice tui.Item, decoded_xsrf string, mavToken string) {
	s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	s.Start()
	animeList, decoded_xsrf, mavToken, err := network.GetAnimeList()
	s.Stop()
	if err != nil {
		panic(err)
	}
	episode_list := []list.Item{}
	for _, anime := range animeList {
		episode_list = append(episode_list, tui.Item{Name: anime[2], Url: anime[1]})
	}
	episode_list = append(episode_list, tui.Item{Name: "Retour", Url: "Retour au menu"})
	m := tui.Model{List: list.New(episode_list, list.NewDefaultDelegate(), 0, 0)}
	m.List.Title = "Tout les animes"

	p := tea.NewProgram(m, tea.WithAltScreen())

	mf, err := p.Run()
	if err != nil {
		panic(err)
	}
	m = mf.(tui.Model)
	if m.Quiting {
		os.Exit(0)
	}
	return m.Choice, decoded_xsrf, mavToken

}
