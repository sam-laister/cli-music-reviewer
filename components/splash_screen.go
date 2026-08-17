package components

import (
	"cli-music-reviewer/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var asciiLogo = `
  ___ _ _      ___       _      _
 / __| (_)_ __| _ ) __ _| |_ __| |_
| (__| | | '_ \ _ \/ _` + "`" + ` |  _/ _| ' \
 \___|_|_| .__/___/\__,_|\__\__|_||_|
         |_|
`

type SplashScreenModel struct {
}

func (s *SplashScreenModel) Update(msg tea.Msg) (*SplashScreenModel, tea.Cmd) {
	return s, nil
}

func (s *SplashScreenModel) View() string {
	subtitle := styles.SubtitleStyle.Render("⚡ Batch TikTok Clip Generator & Auto-Captioner ⚡")

	titlePart := styles.AsciiTitleStyle.Render(strings.TrimSpace(asciiLogo))
	content := lipgloss.JoinVertical(lipgloss.Center, titlePart, subtitle)
	return styles.AsciiHeroStyle.Render(content)
}

func NewSplashScreen() *SplashScreenModel {
	return &SplashScreenModel{}
}
