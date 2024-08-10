package tui

import (
	"fmt"
	"io"
	"strings"

	stplatform "github.com/StevenSermeus/anime-tui/st_platform"
	videoplayer "github.com/StevenSermeus/anime-tui/video_player/vlc"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const listHeight = 20

var docStyle = lipgloss.NewStyle().Margin(1, 2)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

type item string

func (i item) FilterValue() string {
	return string(i)
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type ListModel struct {
	list      list.Model
	Choice    string
	Quitting  bool
	Returning bool
}

func (m ListModel) Init() tea.Cmd {
	return nil
}

func (m ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				if string(i) == "Retour" {
					m.Returning = true
					return m, tea.Quit
				}
				m.Choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ListModel) View() string {
	return docStyle.Render(m.list.View())
}

func StyleList(list *list.Model) {
	list.SetShowStatusBar(true)
	list.SetFilteringEnabled(true)
	list.Styles.Title = titleStyle
	list.Styles.PaginationStyle = paginationStyle
	list.Styles.HelpStyle = helpStyle
}

func Menu() (string, bool) {
	items := []list.Item{
		item("Nouveautés"),
		item("Recherche"),
	}

	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Anime TUI"
	StyleList(&l)

	m := ListModel{list: l}

	end_m, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		return "", true
	}
	end_model := end_m.(ListModel)
	return end_model.Choice, end_model.Quitting
}

func Search(mav stplatform.StreamingPlatform) (stplatform.Anime, bool, bool, error) {
	animes, err := mav.GetAnimeList()
	if err != nil {
		fmt.Println("Error getting anime list:", err)
	}
	items := make([]list.Item, len(animes)+1)
	for i, anime := range animes {
		items[i] = item(anime.Title)
	}
	items[len(animes)] = item("Retour")
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Anime TUI - Recherche"
	StyleList(&l)

	m := ListModel{list: l}

	end_m, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		return stplatform.Anime{}, false, true, err
	}
	end_model := end_m.(ListModel)
	if end_model.Choice == "" {
		return stplatform.Anime{}, end_model.Returning, end_model.Quitting, nil
	}
	for _, anime := range animes {
		if anime.Title == end_model.Choice {
			return anime, end_model.Returning, end_model.Quitting, nil
		}
	}
	return stplatform.Anime{}, end_model.Returning, end_model.Quitting, nil
}

func Newest(mav stplatform.StreamingPlatform) (bool, bool, error) {
	episodes, err := mav.GetRecentEpisode()
	if err != nil {
		return false, true, err
	}
	items := make([]list.Item, len(episodes)+1)
	for i, episode := range episodes {
		items[i] = item(episode.Title)
	}
	items[len(episodes)] = item("Retour")
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Anime TUI - Nouveautés"
	StyleList(&l)

	m := ListModel{list: l}

	end_m, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		return false, true, err
	}
	end_model := end_m.(ListModel)

	if end_model.Choice == "" {
		return end_model.Returning, end_model.Quitting, nil
	}
	anime_choice := stplatform.Anime{}
	for _, episode := range episodes {
		if episode.Title == end_model.Choice {
			anime_choice.Title = episode.Title
			anime_choice.URL = episode.URL
			break
		}
	}
	if anime_choice.Title == "" {
		return end_model.Returning, end_model.Quitting, nil
	}
	video_providers, err := mav.GetVideoURL(anime_choice.URL)
	if err != nil {
		fmt.Println("Error getting video URL:", err)
		return false, true, err
	}
	if len(video_providers) == 0 {
		fmt.Println("No video provider found")
		return false, true, nil
	}
	video_provider := video_providers[0]
	vid_url, err := video_provider.GetVideoUrl()
	if err != nil {
		fmt.Println("Error getting video URL:", err)
		return false, true, err
	}
	vlc := videoplayer.VLC{VideoUrl: vid_url.Url, Title: anime_choice.Title, Refferer: vid_url.Reffer}
	err = vlc.PlayAnime()
	if err != nil {
		fmt.Println("Error playing video:", err)
		return false, true, err
	}
	return false, false, nil
}

func Episodes(mav stplatform.StreamingPlatform, anime stplatform.Anime) (bool, bool, error) {
	episodes, err := mav.GetEpisodeList(anime.URL)
	if err != nil {
		return false, true, err
	}
	items := make([]list.Item, len(episodes)+1)
	for i, episode := range episodes {
		items[i] = item(episode.Title)
	}
	items[len(episodes)] = item("Retour")
	const defaultWidth = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Anime TUI - Episodes"
	StyleList(&l)

	m := ListModel{list: l}

	end_m, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		return false, true, err
	}
	end_model := end_m.(ListModel)

	if end_model.Choice == "" {
		return end_model.Returning, end_model.Quitting, nil
	}
	anime_choice := stplatform.Episode{}
	for _, episode := range episodes {
		if episode.Title == end_model.Choice {
			anime_choice.Title = episode.Title
			anime_choice.URL = episode.URL
			break
		}
	}
	if anime_choice.Title == "" {
		return end_model.Returning, end_model.Quitting, nil
	}
	video_providers, err := mav.GetVideoURL(anime_choice.URL)
	if err != nil {
		fmt.Println("Error getting video URL:", err)
		return false, true, err
	}
	if len(video_providers) == 0 {
		fmt.Println("No video provider found")
		return false, true, nil
	}
	video_provider := video_providers[0]
	vid_url, err := video_provider.GetVideoUrl()
	if err != nil {
		fmt.Println("Error getting video URL:", err)
		return false, true, err
	}
	vlc := videoplayer.VLC{VideoUrl: vid_url.Url, Title: anime_choice.Title, Refferer: vid_url.Reffer}
	err = vlc.PlayAnime()
	if err != nil {
		fmt.Println("Error playing video:", err)
		return false, true, err
	}
	return false, false, nil
}
