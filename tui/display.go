package tui

import (
	"fmt"

	"github.com/jermartinz/hn/api"
	"github.com/jermartinz/hn/models"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/browser"
)

var (
	titleStyle = lipgloss.NewStyle().MarginLeft(2)
	itemStyle  = lipgloss.NewStyle().PaddingLeft(4)
)

type Model struct {
	Items        []*models.Item
	Loading      bool
	Err          error
	Client       *api.APIClient
	Cursor       int
	Page         int
	ItemsPerPage int
	Choice       string
	Quitting     bool
}

type storiesLoadedMsg struct {
	items []*models.Item
}

type errMsg struct {
	err error
}

func (m Model) Init() tea.Cmd {
	return func() tea.Msg {
		items, err := m.Client.GetItemStories()
		if err != nil {
			return errMsg{err}
		}
		return storiesLoadedMsg{items}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case storiesLoadedMsg:
		m.Items = msg.items
		m.Loading = false
		return m, nil
	case errMsg:
		m.Err = msg.err
		m.Loading = false
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			m.Quitting = true
			return m, tea.Quit
		case "enter":
			index := m.Page*m.ItemsPerPage + m.Cursor
			if index >= len(m.Items) {
				return m, nil
			}
			browser.OpenURL(m.Items[index].URL)
			return m, nil
		case "down":
			m.Cursor++
			index := m.Page*m.ItemsPerPage + m.Cursor
			if index >= len(m.Items) {
				m.Cursor--
				return m, nil
			}
			if m.Cursor >= m.ItemsPerPage {
				m.Page++
				m.Cursor = 0
			}
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			} else if m.Page > 0 {
				m.Page--
				m.Cursor = m.ItemsPerPage - 1
			}
		case "right":
			totalPages := len(m.Items) / m.ItemsPerPage
			if len(m.Items)%m.ItemsPerPage > 0 {
				totalPages++
				if m.Page < totalPages-1 {
					m.Page++
				}
			}
		case "left":
			if m.Page > 0 {
				m.Page--
			}
		}
	}
	return m, nil
}

func (m Model) View() string {
	if m.Loading {
		return "Loading Hacker News..."
	}
	if m.Err != nil {
		return fmt.Sprintf("Error: %v", m.Err)
	}
	if len(m.Items) == 0 {
		return "Oops, there are no stories at this time."
	}
	start := m.Page * m.ItemsPerPage
	end := start + m.ItemsPerPage

	if end > len(m.Items) {
		end = len(m.Items)
	}
	var s string
	for i := start; i < end; i++ {
		cursor := " "
		if i-start == m.Cursor {
			cursor = ">"
		}
		s += fmt.Sprintf("%s %d. %s\n", cursor, i+1, m.Items[i].Title)
	}
	return s
}
