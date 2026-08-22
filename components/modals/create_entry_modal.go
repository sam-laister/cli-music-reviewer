package modals

import tea "github.com/charmbracelet/bubbletea"

type CreateEntryModalModel struct{}

func (m *CreateEntryModalModel) Update(msg tea.Msg) (*CreateEntryModalModel, tea.Cmd) {
	return m, nil
}

func (m *CreateEntryModalModel) View() string {}
