package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type Item struct {
	Name, Url string
}

func (i Item) Title() string       { return i.Name }
func (i Item) Description() string { return i.Url }
func (i Item) FilterValue() string { return i.Name }

type Model struct {
	List    list.Model
	Quiting bool
	Choice  Item
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			m.Quiting = true
			return m, tea.Quit
		}
		if msg.String() == "enter" {
			m.Quiting = false
			i, ok := m.List.SelectedItem().(Item)
			if ok {
				m.Choice = i
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)

	}
	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return docStyle.Render(m.List.View())
}
