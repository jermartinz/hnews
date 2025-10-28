package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jermartinz/hn/api"
	"github.com/jermartinz/hn/tui"
)

func main() {
	client := api.NewClient()
	m := tui.Model{
		Client:       client,
		ItemsPerPage: 10,
		Loading:      true,
	}
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
