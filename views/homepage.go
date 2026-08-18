package views

import (
	"cli-music-reviewer/components"
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type HomepageModel struct {
	state       homepageState
	splashPage  *components.SplashScreenModel
	browserPage *components.EntryBrowserModel
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.state = (m.state + 1) % StateCount
		}
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

	instructions := styles.InstructionStyle.Render("Press 'tab' to switch views • 'q' to quit")

	return fmt.Sprintf("\n%s\n\n%s\n", currentView, instructions)
}

func NewHomepage(repos *repositories.AppRepositories) tea.Model {
	return HomepageModel{
		state:       StateSplash,
		splashPage:  components.NewSplashScreen(),
		browserPage: components.NewEntryBrowser(true, repos),
		repos:       repos,
	}
}
