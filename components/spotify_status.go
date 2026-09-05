package components

import (
	"cli-music-reviewer/services"
	"cli-music-reviewer/styles"
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// spotifyAuthCheckedMsg carries the result of an EnsureAuthorized call made
// in response to a refresh request.
type spotifyAuthCheckedMsg struct {
	err error
}

type SpotifyStatusModel struct {
	spotifyHandler services.SpotifyHandler
	authorized     bool
	loading        bool
	spinner        spinner.Model
}

func NewSpotifyStatus(spotifyHandler services.SpotifyHandler) *SpotifyStatusModel {
	return &SpotifyStatusModel{
		spotifyHandler: spotifyHandler,
		authorized:     spotifyHandler.EnsureAuthorized() == nil,
		spinner:        spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
}

func (m *SpotifyStatusModel) checkAuth() tea.Msg {
	return spotifyAuthCheckedMsg{err: m.spotifyHandler.EnsureAuthorized()}
}

func (m *SpotifyStatusModel) Update(msg tea.Msg) (*SpotifyStatusModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "r" && !m.loading {
			m.loading = true
			return m, tea.Batch(m.spinner.Tick, m.checkAuth)
		}
	case spotifyAuthCheckedMsg:
		m.loading = false
		m.authorized = msg.err == nil
	case spinner.TickMsg:
		if !m.loading {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *SpotifyStatusModel) View() string {
	header := styles.ConfigHeaderStyle.Render(" Spotify ")

	var indicator string
	switch {
	case m.loading:
		indicator = m.spinner.View()
	case m.authorized:
		indicator = styles.SuccessStyle.Render("✓")
	default:
		indicator = styles.ErrorStyle.Render("✗")
	}

	row := fmt.Sprintf("  Authorization %s", indicator)
	instructions := styles.InstructionStyle.Render("'r' to refresh")

	content := header + "\n\n" + row + "\n\n" + instructions
	return styles.ConfigHeroStyle.Render(content)
}
