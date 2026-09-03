package views

import (
	"cli-music-reviewer/components"
	"cli-music-reviewer/components/modals"
	"cli-music-reviewer/events"
	"cli-music-reviewer/models/entities"
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/styles"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type HomepageModel struct {
	state       homepageState
	splashPage  *components.SplashScreenModel
	browserPage *components.EntryBrowserModel
	modal       *modals.CreateEntryModalModel
	repos       *repositories.AppRepositories
}

type homepageState int

const (
	StateSplash homepageState = iota
	StateMenu
	StateCount
)

func (m HomepageModel) Init() tea.Cmd {
	return nil
}

func (m HomepageModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	if m.modal != nil {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		case events.EntryCreateSubmittedMsg:
			m.createEntry(msg.Title)
			m.modal = nil
			return m, nil
		case events.EntryCreateCancelledMsg:
			m.modal = nil
			return m, nil
		}

		m.modal, cmd = m.modal.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.state = (m.state + 1) % StateCount
		}
	case events.EntryCreateRequestedMsg:
		var focusCmd tea.Cmd
		m.modal, focusCmd = modals.NewCreateEntryModal()
		return m, focusCmd
	}

	switch m.state {
	case StateSplash:
		m.splashPage, cmd = m.splashPage.Update(msg)
	case StateMenu:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch keypress := msg.String(); keypress {
			case "up":
				m.browserPage.CursorUp()
			case "down":
				m.browserPage.CursorDown()
			}
		}

		m.browserPage, cmd = m.browserPage.Update(msg)
	default:
		panic("unknown state")
	}

	return m, cmd
}

func (m HomepageModel) View() string {
	var currentView string

	switch m.state {
	case StateSplash:
		currentView = m.splashPage.View()
	case StateMenu:
		currentView = m.browserPage.View()
	default:
		panic("unknown state")
	}

	if m.modal != nil {
		currentView = m.modal.View()
	}

	instructions := styles.InstructionStyle.Render("Press 'tab' to switch views • 'q' to quit")

	return fmt.Sprintf("\n%s\n\n%s\n", currentView, instructions)
}

func (m *HomepageModel) createEntry(title string) {
	title = strings.TrimSpace(title)
	if title == "" {
		return
	}

	now := time.Now().UTC().Format(time.RFC3339)
	if _, err := m.repos.EntryRowRepository.Create(&entities.EntryRow{Title: title, CreatedAt: now, UpdatedAt: now, Active: true}); err != nil {
		return
	}

	m.browserPage = components.NewEntryBrowser(true, m.repos)
}

func NewHomepage(repos *repositories.AppRepositories) tea.Model {
	return HomepageModel{
		state:       StateSplash,
		splashPage:  components.NewSplashScreen(),
		browserPage: components.NewEntryBrowser(true, repos),
		repos:       repos,
	}
}
