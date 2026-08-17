package main

import (
	"cli-music-reviewer/views"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(views.NewHomepage()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "Uh oh:", err)
		os.Exit(1)
	}
}
