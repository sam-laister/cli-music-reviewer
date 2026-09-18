package components

import (
	"cli-music-reviewer/config"
	"cli-music-reviewer/events"
	"cli-music-reviewer/interfaces"
	"cli-music-reviewer/models/entities"
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/services"
	"cli-music-reviewer/styles"
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	infoArtCols   = 24
	infoArtRows   = 12
	infoPanelCols = 30
)

type entryArtworkLoadedMsg struct {
	requestID int
	url       string
	rendered  string
	err       error
}

type EntryBrowserModel struct {
	activeIndex  int
	children     []*EntryRowModel
	showControls bool

	artworkService services.ArtworkService
	artCache       map[string]string
	artErr         error
	requestID      int
}

func (m *EntryBrowserModel) Update(msg tea.Msg) (*EntryBrowserModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "c":
			return m, func() tea.Msg { return events.EntryCreateRequestedMsg{} }
		case "enter":
			if len(m.children) == 0 {
				return m, nil
			}
			entryID := m.children[m.activeIndex].id
			return m, func() tea.Msg { return events.EntryEditRequestedMsg{EntryID: entryID} }
		}
	case entryArtworkLoadedMsg:
		if msg.requestID != m.requestID {
			return m, nil
		}
		m.artErr = msg.err
		if msg.err == nil {
			m.artCache[msg.url] = msg.rendered
		}
		return m, nil
	}

	for i := range m.children {
		m.children[i].SetSelected(i == m.activeIndex)
		m.children[i], cmd = m.children[i].Update(msg)
	}
	return m, cmd
}

func (m *EntryBrowserModel) View() string {
	header := styles.ConfigHeaderStyle.Render(fmt.Sprintf(" %s ", config.AppTitle))
	versionPart := lipgloss.NewStyle().Foreground(styles.Gray).Render(" Configuration ")
	headerRow := lipgloss.JoinHorizontal(lipgloss.Bottom, header, versionPart)

	var rows []string
	rows = append(rows, headerRow, "")

	for _, entry := range m.children {
		rows = append(rows, entry.View())
	}

	content := strings.Join(rows, "\n")
	listPane := styles.ConfigHeroStyle.Render(content)

	info := m.infoPanelView()
	if info == "" {
		return listPane
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, listPane, "  ", info)
}

func (m *EntryBrowserModel) infoPanelView() string {
	if len(m.children) == 0 {
		return ""
	}

	entry := m.children[m.activeIndex].entry
	if entry == nil {
		return ""
	}

	artBox := lipgloss.NewStyle().Width(infoArtCols).Height(infoArtRows)

	url := entryArtworkURL(entry)
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

	if entry.Artist != "" {
		lines = append(lines, styles.InstructionStyle.Render("Artist"), entry.Artist, "")
	}
	if entry.ReleaseDate != "" {
		lines = append(lines, styles.InstructionStyle.Render("Released"), entry.ReleaseDate, "")
	}
	if entry.SpotifyLink != "" {
		lines = append(lines, styles.InstructionStyle.Render("Spotify"), entry.SpotifyLink)
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

func (m *EntryBrowserModel) selectedArtworkURL() string {
	if len(m.children) == 0 {
		return ""
	}
	entry := m.children[m.activeIndex].entry
	if entry == nil {
		return ""
	}
	return entryArtworkURL(entry)
}

func (m *EntryBrowserModel) loadArtwork(requestID int, url string) tea.Cmd {
	return func() tea.Msg {
		rendered, err := m.artworkService.Render(url, infoArtCols, infoArtRows)
		return entryArtworkLoadedMsg{requestID: requestID, url: url, rendered: rendered, err: err}
	}
}

func (m *EntryBrowserModel) requestArtwork() tea.Cmd {
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

func (m *EntryBrowserModel) CursorUp() tea.Cmd {
	if m.activeIndex <= 0 {
		m.activeIndex = 0
		return nil
	}
	m.activeIndex--
	return m.requestArtwork()
}

func (m *EntryBrowserModel) CursorDown() tea.Cmd {
	if m.activeIndex >= len(m.children)-1 {
		m.activeIndex = len(m.children) - 1
		return nil
	}
	m.activeIndex++
	return m.requestArtwork()
}

func NewEntryBrowser(showControls bool, repos *repositories.AppRepositories, artworkService services.ArtworkService) (*EntryBrowserModel, tea.Cmd) {
	entryRows, err := repos.EntryRowRepository.GetActiveRows()
	if err != nil {
		log.Printf("entry browser: failed to load active rows: %v", err)
		entryRows = nil
	}

	var rowModels []*EntryRowModel
	for _, row := range entryRows {
		rowModels = append(rowModels, &EntryRowModel{
			isSelected: false,
			id:         row.ID,
			title:      row.Title,
			timestamp:  services.DateToString(row.UpdatedAt),
			entry:      row,
		})
	}

	m := &EntryBrowserModel{
		activeIndex:    0,
		children:       rowModels,
		showControls:   showControls,
		artworkService: artworkService,
		artCache:       make(map[string]string),
	}

	return m, m.requestArtwork()
}

var (
	_ interfaces.ComponentInterface = (*EntryBrowserModel)(nil)
)
