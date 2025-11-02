package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/jermartinz/hnews/tui"
)

func main() {
	m := tui.ModelStyle()

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
