package createentry

import (
	"cli-music-reviewer/components/modals/album_picker"
	"cli-music-reviewer/components/modals/entryform"
	"cli-music-reviewer/services"

	tea "github.com/charmbracelet/bubbletea"
)

type createEntryStep int

const (
	stepPicker createEntryStep = iota
	stepForm
)

type Model struct {
	step   createEntryStep
	picker *album_picker.Model
	form   *entryform.Model

	spotifyHandler services.SpotifyHandler
	artworkService services.ArtworkService
}

func New(spotifyHandler services.SpotifyHandler, artworkService services.ArtworkService) (*Model, tea.Cmd) {
	picker, loadCmd := album_picker.New(spotifyHandler, artworkService)

	m := &Model{
		step:           stepPicker,
		picker:         picker,
		spotifyHandler: spotifyHandler,
		artworkService: artworkService,
	}

	return m, loadCmd
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	switch msg := msg.(type) {
	case album_picker.AlbumSelectedMsg:
		artwork := m.picker.GetCacheItem(services.AlbumArtworkURL(msg.Album))
		m.form = entryform.New(msg.Album, artwork)
		m.step = stepForm
		return m, nil
	case entryform.BackToPickerMsg:
		m.step = stepPicker
		m.form = nil
		return m, nil
	}

	var cmd tea.Cmd
	switch m.step {
	case stepPicker:
		m.picker, cmd = m.picker.Update(msg)
	case stepForm:
		m.form, cmd = m.form.Update(msg)
	}
	return m, cmd
}

func (m *Model) View() string {
	switch m.step {
	case stepForm:
		return m.form.View()
	default:
		return m.picker.View()
	}
}
