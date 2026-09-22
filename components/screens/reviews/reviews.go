package reviews

import (
	"cli-music-reviewer/models/entities"
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/services"
	"cli-music-reviewer/styles"
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sam-laister/sam-laister-bubbletea-components-library/list"
)

const (
	infoArtCols       = 24
	infoArtRows       = 12
	infoPanelCols     = 30
	reviewsListWidth  = 50
	reviewsListHeight = 10
)

type Model struct {
	list           list.Model[*entities.EntryRow]
	showControls   bool
	artworkService services.ArtworkService
	artCache       map[string]string
	artErr         error
	requestID      int
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	var artCmd tea.Cmd

	switch msg := msg.(type) {
	case list.CursorMsg[*entities.EntryRow]:
		artCmd = m.requestArtwork()
	case entryArtworkLoadedMsg:
		if msg.requestID != m.requestID {
			return m, nil
		}
		m.artErr = msg.err
		if msg.err == nil {
			m.artCache[msg.url] = msg.rendered
		}
	case tea.WindowSizeMsg:
		m.list.SetSize(reviewsListWidth, reviewsListHeight)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, tea.Batch(artCmd, cmd)
}

func (m *Model) View() string {
	header := styles.ConfigHeaderStyle.Render(" Reviews ")

	content := fmt.Sprintf("%s\n\n %s\n", header, m.list.View())
	listPane := styles.ConfigHeroStyle.Render(content)

	info := m.infoPanelView()
	if info == "" {
		return listPane
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, listPane, "  ", info)
}

func (m *Model) infoPanelView() string {
	if m.list.IsEmpty() || m.list.SelectedItem() == nil {
		return ""
	}

	artBox := lipgloss.NewStyle().Width(infoArtCols).Height(infoArtRows)

	url := entryArtworkURL(m.list.SelectedItem())
	var art string
	switch {
	case url == "":
		art = styles.InstructionStyle.Render("[no cover art]")
	case m.artCache[url] != "":
		art = m.artCache[url]
	case m.artErr != nil:
		art = styles.InstructionStyle.Render("[cover art unavailable]")
	default:
		art = styles.InstructionStyle.Render("loading art…")
	}

	lines := []string{artBox.Render(art), ""}

	if m.list.SelectedItem().Artist != "" {
		lines = append(lines, styles.InstructionStyle.Render("Artist"), m.list.SelectedItem().Artist, "")
	}
	if m.list.SelectedItem().ReleaseDate != "" {
		lines = append(lines, styles.InstructionStyle.Render("Released"), m.list.SelectedItem().ReleaseDate, "")
	}
	if m.list.SelectedItem().SpotifyLink != "" {
		lines = append(lines, styles.InstructionStyle.Render("Spotify"), m.list.SelectedItem().SpotifyLink)
	}

	content := strings.Join(lines, "\n")
	return styles.ConfigHeroStyle.Width(infoPanelCols).Render(content)
}

func entryArtworkURL(entry *entities.EntryRow) string {
	if entry.CoverArtMedium != "" {
		return entry.CoverArtMedium
	}
	return entry.CoverArtLarge
}

func (m *Model) selectedArtworkURL() string {
	if m.list.IsEmpty() {
		return ""
	}
	if m.list.SelectedItem() == nil {
		return ""
	}
	return entryArtworkURL(m.list.SelectedItem())
}

func (m *Model) loadArtwork(requestID int, url string) tea.Cmd {
	return func() tea.Msg {
		rendered, err := m.artworkService.Render(url, infoArtCols, infoArtRows)
		return entryArtworkLoadedMsg{requestID: requestID, url: url, rendered: rendered, err: err}
	}
}

func (m *Model) requestArtwork() tea.Cmd {
	if m.artworkService == nil {
		return nil
	}

	url := m.selectedArtworkURL()
	if url == "" {
		return nil
	}
	if _, ok := m.artCache[url]; ok {
		return nil
	}

	m.artErr = nil
	m.requestID++
	return m.loadArtwork(m.requestID, url)
}

func New(showControls bool, repos *repositories.AppRepositories, artworkService services.ArtworkService) (*Model, tea.Cmd) {
	reviews, err := repos.EntryRowRepository.GetActiveRows()
	if err != nil {
		log.Printf("entry browser: failed to load active rows: %v", err)
		reviews = nil
	}

	m := &Model{
		list:           list.New(reviews, reviewDelegate{}),
		showControls:   showControls,
		artworkService: artworkService,
		artCache:       make(map[string]string),
	}

	m.list.SetSize(reviewsListWidth, reviewsListHeight)

	return m, m.requestArtwork()
}
