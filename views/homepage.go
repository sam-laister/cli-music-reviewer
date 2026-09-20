package views

import (
	"cli-music-reviewer/components"
	"cli-music-reviewer/components/modals"
	"cli-music-reviewer/events"
	"cli-music-reviewer/models/entities"
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/services"
	"cli-music-reviewer/styles"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HomepageModel struct {
	state             homepageState
	splashPage        *components.SplashScreenModel
	browserPage       *components.EntryBrowserModel
	spotifyStatusPage *components.SpotifyStatusModel
	modal             *modals.CreateEntryModalModel
	reviewEditor      *modals.ReviewEditorModel
	repos             *repositories.AppRepositories
	services          *services.AppServices
	termWidth         int
	termHeight        int
	initialCmd        tea.Cmd
}

type homepageState int

const (
	StateSplash homepageState = iota
	StateMenu
	StateSpotifyStatus
	StateCount
)

func (m HomepageModel) Init() tea.Cmd {
	return m.initialCmd
}

func (m HomepageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		m.termWidth, m.termHeight = sizeMsg.Width, sizeMsg.Height
	}

	if m.modal != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		case events.EntryCreateSubmittedMsg:
			browserCmd := m.createEntry(msg)
			m.modal = nil
			return m, tea.Batch(browserCmd, tea.ClearScreen)
		case events.EntryCreateCancelledMsg:
			m.modal = nil
			return m, tea.ClearScreen
		}

		m.modal, cmd = m.modal.Update(msg)
		return m, cmd
	}

	if m.reviewEditor != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		case events.ReviewSaveRequestedMsg:
			_ = m.repos.EntryRowRepository.Update(msg.Entry)
			m.reviewEditor = nil
			var browserCmd tea.Cmd
			m.browserPage, browserCmd = components.NewEntryBrowser(true, m.repos, m.services.ArtworkService)
			return m, tea.Batch(browserCmd, tea.ClearScreen)
		case events.ReviewEditCancelledMsg:
			m.reviewEditor = nil
			return m, tea.ClearScreen
		}

		m.reviewEditor, cmd = m.reviewEditor.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.state = (m.state + 1) % StateCount
			return m, tea.ClearScreen
		}
	case events.EntryCreateRequestedMsg:
		var focusCmd tea.Cmd
		m.modal, focusCmd = modals.NewCreateEntryModal(m.services.SpotifyHandler, m.services.ArtworkService)
		return m, tea.Batch(focusCmd, tea.ClearScreen)
	case events.EntryEditRequestedMsg:
		if editor, err := modals.NewReviewEditor(msg.EntryID, m.repos.EntryRowRepository, m.termWidth, m.termHeight); err == nil {
			m.reviewEditor = editor
			return m, tea.ClearScreen
		}
		return m, nil
	}

	switch m.state {
	case StateSplash:
		m.splashPage, cmd = m.splashPage.Update(msg)
	case StateMenu:
		var navCmd tea.Cmd
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "up":
				navCmd = m.browserPage.CursorUp()
			case "down":
				navCmd = m.browserPage.CursorDown()
			}
		}

		var updateCmd tea.Cmd
		m.browserPage, updateCmd = m.browserPage.Update(msg)
		cmd = tea.Batch(navCmd, updateCmd)
	case StateSpotifyStatus:
		m.spotifyStatusPage, cmd = m.spotifyStatusPage.Update(msg)
	default:
		panic("unknown state")
	}

	return m, cmd
}

func (m HomepageModel) View() string {
	// The editor is a distinct full-screen mode with its own neovim-style
	// (ctrl+key) keybindings, not part of the tab-cycling "central modal"
	// paradigm the rest of the app uses — it gets the whole screen, nothing
	// centered around it, and no generic tab/quit instructions line (tab and
	// q are captured by the editor itself while it's focused).
	if m.reviewEditor != nil {
		return m.services.ArtworkService.ClearImages() + m.reviewEditor.View()
	}

	var currentView string
	switch m.state {
	case StateSplash:
		currentView = m.splashPage.View()
	case StateMenu:
		currentView = m.browserPage.View()
	case StateSpotifyStatus:
		currentView = m.spotifyStatusPage.View()
	default:
		panic("unknown state")
	}

	canShowArtwork := m.state == StateMenu
	if m.modal != nil {
		currentView = m.modal.View()
		canShowArtwork = true
	}

	// See services/artwork_service.go / album picker & create-entry modal —
	// those views bake their own Kitty delete-then-place into whatever
	// cached image string they re-embed, so nothing extra is needed when
	// they're on screen. When neither can be showing artwork, prepend a
	// clear so nothing stale from a previous frame lingers. This must stay
	// conditional (present only some frames), not unconditional — browser/
	// modal views share a leading blank margin line (styles.ConfigHeroStyle
	// / styles.ModalStyle both use Margin(1,0,1,0)), so an always-on prefix
	// would sit on an always-identical line and Bubble Tea's line-diffing
	// renderer would permanently skip resending it after the first paint.
	// Toggling the prefix's presence on/off exactly at these transitions is
	// what guarantees the line actually changes when it needs to.
	if !canShowArtwork {
		currentView = m.services.ArtworkService.ClearImages() + currentView
	}

	instructions := styles.InstructionStyle.Render("Press 'tab' to switch views • 'q' to quit")
	block := currentView + "\n\n" + instructions

	if m.termWidth <= 0 || m.termHeight <= 0 {
		// No WindowSizeMsg yet (first frame) — render un-centered rather
		// than collapse into a 0x0 canvas.
		return "\n" + block + "\n"
	}

	return lipgloss.Place(m.termWidth, m.termHeight, lipgloss.Center, lipgloss.Center, block)
}

func (m *HomepageModel) createEntry(msg events.EntryCreateSubmittedMsg) tea.Cmd {
	title := strings.TrimSpace(msg.Title)
	if title == "" {
		return nil
	}

	entry := entities.NewEntryRow(
		title,
		"",
		msg.Artist,
		msg.ReleaseDate,
		msg.SpotifyID,
		msg.SpotifyType,
		msg.SpotifyLink,
		msg.CoverArtSmall,
		msg.CoverArtMedium,
		msg.CoverArtLarge,
		true,
	)

	if _, err := m.repos.EntryRowRepository.Create(entry); err != nil {
		return nil
	}

	var browserCmd tea.Cmd
	m.browserPage, browserCmd = components.NewEntryBrowser(true, m.repos, m.services.ArtworkService)
	return browserCmd
}

func NewHomepage(repos *repositories.AppRepositories, services *services.AppServices) tea.Model {
	browserPage, browserCmd := components.NewEntryBrowser(true, repos, services.ArtworkService)

	return HomepageModel{
		state:             StateSplash,
		splashPage:        components.NewSplashScreen(),
		browserPage:       browserPage,
		spotifyStatusPage: components.NewSpotifyStatus(services.SpotifyHandler, services.HttpHandler),
		repos:             repos,
		services:          services,
		initialCmd:        browserCmd,
	}
}
