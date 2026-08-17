package interfaces

import tea "github.com/charmbracelet/bubbletea"

type ComponentInterface interface {
	Update(msg tea.Msg) (*ComponentInterface, tea.Cmd)
	View() string
}
