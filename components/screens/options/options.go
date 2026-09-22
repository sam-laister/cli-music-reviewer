package options

import (
	"cli-music-reviewer/services"
	"cli-music-reviewer/styles"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sam-laister/sam-laister-bubbletea-components-library/checkbox"
	optionslist "github.com/sam-laister/sam-laister-bubbletea-components-library/options-list"
	textinput "github.com/sam-laister/sam-laister-bubbletea-components-library/text-input"
	"github.com/sam-laister/sam-laister-bubbletea-components-library/theme"
)

const (
	optionsListWidth  = 50
	optionsListHeight = 10
)

type Model struct {
	spotifyHandler services.SpotifyHandler
	httpHandler    services.HttpHandler
	authorized     bool
	loading        bool
	connecting     bool
	err            error
	spinner        spinner.Model
	list           optionslist.Model
}

func (m *Model) View() string {
	header := styles.ConfigHeaderStyle.Render(" Options ")
	return styles.ConfigHeroStyle.Render(header + "\n\n" + m.list.View())
}

func (m *Model) SetSize(width, height int) {
	m.list.SetSize(width, height)
}

func (m *Model) checkAuth() tea.Msg {
	return spotifyAuthCheckedMsg{err: m.spotifyHandler.EnsureAuthorized()}
}

// connect runs the full browser Authorize/callback flow, bypassing
// EnsureAuthorized's reuse-if-present shortcut — this is what makes 'a' a
// genuine re-authorize action even if a (possibly stale-scoped) token is
// already stored.
func (m *Model) connect() tea.Msg {
	if err := m.httpHandler.Setup(); err != nil {
		return spotifyConnectedMsg{err: err}
	}
	if err := m.spotifyHandler.Authorize(); err != nil {
		return spotifyConnectedMsg{err: err}
	}
	return spotifyConnectedMsg{err: m.httpHandler.Wait()}
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			if !m.loading && !m.connecting && !m.list.Focused() {
				m.loading = true
				return m, tea.Batch(m.spinner.Tick, m.checkAuth)
			}
		case "a":
			if !m.authorized && !m.loading && !m.connecting && !m.list.Focused() {
				m.connecting = true
				m.err = nil
				return m, tea.Batch(m.spinner.Tick, m.connect)
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(optionsListWidth, optionsListHeight)
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
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func New(spotifyHandler services.SpotifyHandler, httpHandler services.HttpHandler) *Model {
	items := []optionslist.Field{
		checkbox.New("Kid A", false, checkbox.WithTheme(theme.Default())),
		checkbox.New("Blonde", false, checkbox.WithTheme(theme.Default())),
		checkbox.New("To Pimp a Butterfly", false, checkbox.WithTheme(theme.Default())),
		textinput.New("In Rainbows", "10", textinput.WithTheme(theme.Default())),
		textinput.New("Channel Orange", "5", textinput.WithTheme(theme.Default())),
	}

	list := optionslist.New(items, theme.Default())
	list.SetSize(optionsListWidth, optionsListHeight)

	return &Model{
		spotifyHandler: spotifyHandler,
		httpHandler:    httpHandler,
		authorized:     spotifyHandler.EnsureAuthorized() == nil,
		spinner:        spinner.New(spinner.WithSpinner(spinner.Dot)),
		list:           list,
	}
}
