package components

import (
	"cli-music-reviewer/components/options"
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

// spotifyConnectedMsg carries the result of a full Authorize/callback flow
// started in response to an authorize request.
type spotifyConnectedMsg struct {
	err error
}

type OptionsScreenModel struct {
	spotifyHandler services.SpotifyHandler
	httpHandler    services.HttpHandler
	authorized     bool
	loading        bool
	connecting     bool
	err            error
	spinner        spinner.Model
	activeIndex    int
	children       []*options.OptionTickboxModel
}

func NewOptionsScreen(spotifyHandler services.SpotifyHandler, httpHandler services.HttpHandler) *OptionsScreenModel {
	return &OptionsScreenModel{
		spotifyHandler: spotifyHandler, 
		httpHandler:    httpHandler,
		authorized:     spotifyHandler.EnsureAuthorized() == nil,
		spinner:        spinner.New(spinner.WithSpinner(spinner.Dot)),
		activeIndex:    0,
	}
}

func (m *OptionsScreenModel) checkAuth() tea.Msg {
	return spotifyAuthCheckedMsg{err: m.spotifyHandler.EnsureAuthorized()}
}

// connect runs the full browser Authorize/callback flow, bypassing
// EnsureAuthorized's reuse-if-present shortcut — this is what makes 'a' a
// genuine re-authorize action even if a (possibly stale-scoped) token is
// already stored.
func (m *OptionsScreenModel) connect() tea.Msg {
	if err := m.httpHandler.Setup(); err != nil {
		return spotifyConnectedMsg{err: err}
	}
	if err := m.spotifyHandler.Authorize(); err != nil {
		return spotifyConnectedMsg{err: err}
	}
	return spotifyConnectedMsg{err: m.httpHandler.Wait()}
}

func (m *OptionsScreenModel) Update(msg tea.Msg) (*OptionsScreenModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if !m.loading && !m.connecting {
				m.loading = true
				return m, tea.Batch(m.spinner.Tick, m.checkAuth)
			}
		case "a":
			if !m.authorized && !m.loading && !m.connecting {
				m.connecting = true
				m.err = nil
				return m, tea.Batch(m.spinner.Tick, m.connect)
			}
		}
	case spotifyAuthCheckedMsg:
		m.loading = false
		m.authorized = msg.err == nil
	case spotifyConnectedMsg:
		m.connecting = false
		m.err = msg.err
		m.authorized = msg.err == nil
	case spinner.TickMsg:
		if !m.loading && !m.connecting {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m *OptionsScreenModel) View() string {
	header := styles.ConfigHeaderStyle.Render(" Options ")

	var indicator string
	switch {
	case m.loading, m.connecting:
		indicator = m.spinner.View()
	case m.authorized:
		indicator = styles.SuccessStyle.Render("✓")
	default:
		indicator = styles.ErrorStyle.Render("✗")
	}

	row := fmt.Sprintf("  Authorization %s", indicator)

	var instructions string
	switch {
	case m.connecting:
		instructions = styles.InstructionStyle.Render("waiting for browser authorization…")
	case !m.authorized:
		instructions = styles.InstructionStyle.Render("'r' to refresh • 'a' to authorize")
	default:
		instructions = styles.InstructionStyle.Render("'r' to refresh")
	}

	content := header + "\n\n" + row
	if m.err != nil {
		content += "\n\n" + styles.ErrorStyle.Render(m.err.Error())
	}
	content += "\n\n" + instructions

	return styles.ConfigHeroStyle.Render(content)
}

func (m *OptionsScreenModel) CursorUp() tea.Cmd {
	if m.activeIndex <= 0 {
		m.activeIndex = 0
		return nil
	}
	m.activeIndex--
	return nil
}

func (m *OptionsScreenModel) CursorDown() tea.Cmd {
	if m.activeIndex >= len(m.children)-1 {
		m.activeIndex = len(m.children) - 1
		return nil
	}
	m.activeIndex++
	return nil
}
