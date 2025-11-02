package tui

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jermartinz/hnews/api"
	"github.com/jermartinz/hnews/models"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/paginator"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	orangeColor = lipgloss.Color("#ff6600")
	appStyle    = lipgloss.NewStyle().Padding(1, 2)
	titleSyte   = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#16", Dark: "255"}).
			Background(orangeColor).
			Padding(0, 1)
)

func ModelStyle() Model {
	items := []list.Item{}
	d := list.NewDefaultDelegate()
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(orangeColor).BorderLeftForeground(orangeColor)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(orangeColor).BorderLeftForeground(orangeColor)
	l := list.New(items, d, 80, 20)
	l.SetShowPagination(true)

	l.Paginator.PerPage = 10
	l.Title = "Hacker News"
	l.Styles.Title = titleSyte
	l.Paginator.Type = paginator.Dots
	l.Paginator.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	l.Paginator.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	return Model{
		List:    l,
		Loading: true,
		Client:  api.NewClient(),
	}
}

type Model struct {
	Items   []*models.Item
	List    list.Model
	Loading bool
	Err     error
	Client  *api.APIClient
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

func isWSL() bool {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return false
	}
	return checkWsl(data)
}

func checkWsl(data []byte) bool {
	return strings.Contains(strings.ToLower(string(data)), "microsoft")
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case storiesLoadedMsg:
		m.Items = msg.items
		items := make([]list.Item, len(msg.items))

		for i, item := range msg.items {
			items[i] = item
		}
		m.List.SetItems(items)
		m.Loading = false
		return m, nil
	case errMsg:
		m.Err = msg.err
		m.Loading = false
		return m, nil
	case tea.WindowSizeMsg:
		m.List.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			selectedItem := m.List.SelectedItem()
			if item, ok := selectedItem.(*models.Item); ok {
				go func() {
					var cmd *exec.Cmd
					switch runtime.GOOS {
					case "linux":
						if isWSL() {
							cmd = exec.Command("wsl-open", item.URL)
						} else {
							cmd = exec.Command("xdg-open", item.URL)
						}
					case "darwin":
						cmd = exec.Command("open", item.URL)
					case "windows":
						cmd = exec.Command("cmd", "/c", "start", item.URL)
					}
					if cmd != nil {
						cmd.Stdout = nil
						cmd.Stderr = nil
						cmd.Start()
					}
				}()
			}
			return m, nil

		}
	}
	m.List, cmd = m.List.Update(msg)
	return m, cmd
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
	return appStyle.Render(m.List.View())
}
