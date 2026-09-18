package modals

import (
	"cli-music-reviewer/services"

	tea "github.com/charmbracelet/bubbletea"
)

type createEntryStep int

const (
	stepPicker createEntryStep = iota
	stepForm
)

type CreateEntryModalModel struct {
	step   createEntryStep
	picker *albumPickerModel
	form   *entryFormModel

	spotifyHandler services.SpotifyHandler
	artworkService services.ArtworkService
}

func NewCreateEntryModal(spotifyHandler services.SpotifyHandler, artworkService services.ArtworkService) (*CreateEntryModalModel, tea.Cmd) {
	picker, loadCmd := newAlbumPicker(spotifyHandler, artworkService)

	m := &CreateEntryModalModel{
		step:           stepPicker,
		picker:         picker,
		spotifyHandler: spotifyHandler,
		artworkService: artworkService,
	}

	return m, tea.Batch(tea.EnterAltScreen, loadCmd)
}

func (m *CreateEntryModalModel) Update(msg tea.Msg) (*CreateEntryModalModel, tea.Cmd) {
	switch msg := msg.(type) {
	case albumSelectedMsg:
		artwork := m.picker.artCache[albumArtworkURL(msg.album)]
		m.form = newEntryForm(msg.album, artwork)
		m.step = stepForm
		return m, nil
	case backToPickerMsg:
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

func (m *CreateEntryModalModel) View() string {
	switch m.step {
	case stepForm:
		return m.form.View()
	default:
		return m.picker.View()
	}
}
