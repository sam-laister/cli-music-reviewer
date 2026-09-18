package modals

import (
	"cli-music-reviewer/events"
	"cli-music-reviewer/models/dtos"
	"cli-music-reviewer/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type backToPickerMsg struct{}

const (
	focusTitle = iota
	focusArtist
	focusDate
	focusFieldCount
)

type entryFormModel struct {
	album   dtos.AlbumDTO
	artwork string

	inputs     [focusFieldCount]textinput.Model
	focusIndex int
}

func newEntryForm(album dtos.AlbumDTO, artwork string) *entryFormModel {
	title := textinput.New()
	title.Placeholder = "Title"
	title.SetValue(album.Name)
	title.Width = 40
	title.Focus()

	artist := textinput.New()
	artist.Placeholder = "Artist"
	artist.SetValue(albumArtists(album))
	artist.Width = 40

	date := textinput.New()
	date.Placeholder = "Date released"
	date.SetValue(album.ReleaseDate)
	date.Width = 40

	return &entryFormModel{
		album:   album,
		artwork: artwork,
		inputs:  [focusFieldCount]textinput.Model{title, artist, date},
	}
}

func (m *entryFormModel) setFocus(index int) {
	m.focusIndex = index
	for i := range m.inputs {
		if i == index {
			m.inputs[i].Focus()
		} else {
			m.inputs[i].Blur()
		}
	}
}

func (m *entryFormModel) submit() tea.Msg {
	small, medium, large := albumCoverArt(m.album)
	return events.EntryCreateSubmittedMsg{
		Title:          strings.TrimSpace(m.inputs[focusTitle].Value()),
		Artist:         strings.TrimSpace(m.inputs[focusArtist].Value()),
		ReleaseDate:    strings.TrimSpace(m.inputs[focusDate].Value()),
		SpotifyID:      m.album.ID,
		SpotifyType:    m.album.Type,
		SpotifyLink:    m.album.ExternalURLs.Spotify,
		CoverArtSmall:  small,
		CoverArtMedium: medium,
		CoverArtLarge:  large,
	}
}

func (m *entryFormModel) Update(msg tea.Msg) (*entryFormModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return backToPickerMsg{} }
		case "enter":
			return m, m.submit
		case "tab":
			m.setFocus((m.focusIndex + 1) % focusFieldCount)
			return m, nil
		case "shift+tab":
			m.setFocus((m.focusIndex - 1 + focusFieldCount) % focusFieldCount)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

func (m *entryFormModel) View() string {
	header := styles.ConfigHeaderStyle.Render(" New entry ")

	labels := []string{"Title", "Artist", "Date released"}
	var fields []string
	for i, label := range labels {
		fields = append(fields, fmt.Sprintf("%s\n%s", label, m.inputs[i].View()))
	}
	form := lipgloss.NewStyle().Width(44).Render(strings.Join(fields, "\n\n"))

	art := m.artwork
	if art == "" {
		art = styles.InstructionStyle.Render("[cover art unavailable]")
	}
	artPanel := lipgloss.NewStyle().Width(artworkCols).Height(artworkRows).Render(art)

	body := lipgloss.JoinHorizontal(lipgloss.Top, form, "  ", artPanel)
	instructions := styles.InstructionStyle.Render("tab to switch fields • enter to save • esc to go back")
	content := header + "\n\n" + body + "\n\n" + instructions
	return styles.ModalStyle.Render(content)
}
