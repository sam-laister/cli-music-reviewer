package modals

import (
	"cli-music-reviewer/events"
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/services"
	"cli-music-reviewer/styles"
	"errors"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	artworkCols     = 24
	artworkRows     = 12
	visibleListRows = 8
	savedAlbumLimit = 10
)

type pickerState int

const (
	pickerLoading pickerState = iota
	pickerError
	pickerLoaded
)

type albumsLoadedMsg struct {
	albums []dtos.SavedAlbumDTO
	err    error
}

type artworkRenderedMsg struct {
	requestID int
	url       string
	rendered  string
	err       error
}

type albumSelectedMsg struct {
	album dtos.AlbumDTO
}

type albumPickerModel struct {
	spotifyHandler services.SpotifyHandler
	artworkService services.ArtworkService

	state pickerState
	err   error

	albums       []dtos.SavedAlbumDTO
	activeIndex  int
	scrollOffset int

	artCache  map[string]string
	artLoaded bool
	artErr    error
	requestID int

	spinner spinner.Model
}

func newAlbumPicker(spotifyHandler services.SpotifyHandler, artworkService services.ArtworkService) (*albumPickerModel, tea.Cmd) {
	m := &albumPickerModel{
		spotifyHandler: spotifyHandler,
		artworkService: artworkService,
		state:          pickerLoading,
		artCache:       make(map[string]string),
		spinner:        spinner.New(spinner.WithSpinner(spinner.Dot)),
	}
	return m, tea.Batch(m.spinner.Tick, m.loadAlbums)
}

func (m *albumPickerModel) loadAlbums() tea.Msg {
	if err := m.spotifyHandler.EnsureAuthorized(); err != nil {
		return albumsLoadedMsg{err: err}
	}

	resp, err := m.spotifyHandler.GetSavedAlbums(savedAlbumLimit, 0, "")
	if err != nil {
		return albumsLoadedMsg{err: err}
	}

	return albumsLoadedMsg{albums: resp.Items}
}

func (m *albumPickerModel) loadArtwork(requestID int, url string) tea.Cmd {
	return func() tea.Msg {
		rendered, err := m.artworkService.Render(url, artworkCols, artworkRows)
		return artworkRenderedMsg{requestID: requestID, url: url, rendered: rendered, err: err}
	}
}

func (m *albumPickerModel) selectedURL() string {
	if m.state != pickerLoaded || len(m.albums) == 0 {
		return ""
	}
	return albumArtworkURL(m.albums[m.activeIndex].Album)
}

func (m *albumPickerModel) requestArtwork() tea.Cmd {
	url := m.selectedURL()
	if url == "" {
		return nil
	}
	if _, ok := m.artCache[url]; ok {
		return nil
	}

	m.requestID++
	return m.loadArtwork(m.requestID, url)
}

func (m *albumPickerModel) Update(msg tea.Msg) (*albumPickerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case albumsLoadedMsg:
		if msg.err != nil {
			m.state = pickerError
			m.err = msg.err
			return m, nil
		}
		if len(msg.albums) == 0 {
			m.state = pickerError
			m.err = errors.New("no saved albums found on your Spotify account")
			return m, nil
		}
		m.state = pickerLoaded
		m.albums = msg.albums
		return m, m.requestArtwork()

	case artworkRenderedMsg:
		if msg.requestID != m.requestID {
			return m, nil
		}
		m.artErr = msg.err
		if msg.err == nil {
			m.artCache[msg.url] = msg.rendered
		}
		return m, nil

	case spinner.TickMsg:
		if m.state != pickerLoading {
			return m, nil
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return events.EntryCreateCancelledMsg{} }
		}

		if m.state != pickerLoaded {
			return m, nil
		}
		switch msg.String() {
		case "up":
			m.cursorUp()
			return m, m.requestArtwork()
		case "down":
			m.cursorDown()
			return m, m.requestArtwork()
		case "enter":
			return m, func() tea.Msg { return albumSelectedMsg{album: m.albums[m.activeIndex].Album} }
		}
	}

	return m, nil
}

func (m *albumPickerModel) cursorUp() {
	if m.activeIndex <= 0 {
		m.activeIndex = 0
		return
	}
	m.activeIndex--
	if m.activeIndex < m.scrollOffset {
		m.scrollOffset = m.activeIndex
	}
}

func (m *albumPickerModel) cursorDown() {
	if m.activeIndex >= len(m.albums)-1 {
		m.activeIndex = len(m.albums) - 1
		return
	}
	m.activeIndex++
	if m.activeIndex >= m.scrollOffset+visibleListRows {
		m.scrollOffset = m.activeIndex - visibleListRows + 1
	}
}

func (m *albumPickerModel) View() string {
	header := styles.ConfigHeaderStyle.Render(" Pick an album ")

	var body string
	switch m.state {
	case pickerLoading:
		body = fmt.Sprintf("%s Loading your saved albums…", m.spinner.View())
	case pickerError:
		body = styles.ErrorStyle.Render(m.err.Error())
		if errors.Is(m.err, services.ErrNoStoredToken) {
			body += "\n\n" + styles.InstructionStyle.Render("Run `go run ./cmd/run_spotify_auth.go` to connect Spotify, then try again.")
		}
	case pickerLoaded:
		body = lipgloss.JoinHorizontal(lipgloss.Top, m.listView(), "  ", m.artworkView())
	}

	instructions := styles.InstructionStyle.Render("↑/↓ to browse • enter to select • esc to cancel")
	content := header + "\n\n" + body + "\n\n" + instructions
	return styles.ModalStyle.Render(content)
}

func (m *albumPickerModel) listView() string {
	end := min(m.scrollOffset+visibleListRows, len(m.albums))

	var rows []string
	for i := m.scrollOffset; i < end; i++ {
		album := m.albums[i].Album
		label := fmt.Sprintf("%s — %s (%s)", album.Name, albumArtists(album), albumYear(album))

		cursor := "  "
		if i == m.activeIndex {
			cursor = styles.CursorStyle.Render("> ")
			label = styles.SelectedItemStyle.Render(label)
		}
		rows = append(rows, cursor+label)
	}

	return lipgloss.NewStyle().Width(50).Render(strings.Join(rows, "\n"))
}

func (m *albumPickerModel) artworkView() string {
	box := lipgloss.NewStyle().Width(artworkCols).Height(artworkRows)

	url := m.selectedURL()
	if rendered, ok := m.artCache[url]; ok {
		return box.Render(rendered)
	}
	if m.artErr != nil {
		return box.Render(styles.InstructionStyle.Render("[cover art unavailable]"))
	}
	return box.Render(styles.InstructionStyle.Render("loading art…"))
}
