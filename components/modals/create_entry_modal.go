package modals

import (
	"cli-music-reviewer/events"
	"cli-music-reviewer/styles"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateEntryModalModel struct {
	titleInput textinput.Model
}

func NewCreateEntryModal() (*CreateEntryModalModel, tea.Cmd) {
	ti := textinput.New()
	ti.Placeholder = "Album or track title"
	ti.CharLimit = 200
	ti.Width = 40
	cmd := ti.Focus()

	return &CreateEntryModalModel{titleInput: ti}, cmd
}

func (m *CreateEntryModalModel) Update(msg tea.Msg) (*CreateEntryModalModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, func() tea.Msg { return events.EntryCreateCancelledMsg{} }
		case "enter":
			title := strings.TrimSpace(m.titleInput.Value())
			return m, func() tea.Msg { return events.EntryCreateSubmittedMsg{Title: title} }
		}
	}

	var cmd tea.Cmd
	m.titleInput, cmd = m.titleInput.Update(msg)
	return m, cmd
}

func (m *CreateEntryModalModel) View() string {
	content := fmt.Sprintf("New Entry\n\nTitle\n%s\n\nenter to save • esc to cancel", m.titleInput.View())
	return styles.ModalStyle.Render(content)
}
