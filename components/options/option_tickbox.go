package options

import (
	"cli-music-reviewer/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type OptionTickboxModel struct {
	id         string
	isSelected bool
	title      string
}

func (m *OptionTickboxModel) Update(msg tea.Msg) (*OptionTickboxModel, tea.Cmd) {
	return m, nil
}

func (m *OptionTickboxModel) View() string {
	cursor := "  "
	label := m.title

	value := "[ %s ]"

	if m.isSelected {
		cursor = styles.CursorStyle.Render("> ")
		label = styles.SelectedItemStyle.Render(label)
		value = styles.SelectedItemStyle.Render(value)
	}

	return fmt.Sprintf("%s%s %s", cursor, label, value)
}

func (m *OptionTickboxModel) SetSelected(isSelected bool) {
	m.isSelected = isSelected
}

func NewOptionTickboxModel(id string, isSelected bool, title string) *OptionTickboxModel {
	return &OptionTickboxModel{
		id:         id,
		title:      title,
		isSelected: isSelected,
	}
}
