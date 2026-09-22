package splash

import (
	"cli-music-reviewer/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var asciiLogo = `
  ____ _   _ _____
 / ___| | | | ____|
| |   | | | |  _|
| |___| |_| | |___
 \____|\___/|_____|
`

type Model struct {
}

func (s *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	return s, nil
}

func (s *Model) View() string {
	subtitle := styles.SubtitleStyle.Render("🎧 Every album, cued and reviewed. 🎧")

	titlePart := styles.AsciiTitleStyle.Render(strings.Trim(asciiLogo, "\n"))
	content := lipgloss.JoinVertical(lipgloss.Center, titlePart, subtitle)
	return styles.AsciiHeroStyle.Render(content)
}

func New() *Model {
	return &Model{}
}
