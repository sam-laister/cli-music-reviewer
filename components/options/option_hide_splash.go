package options

import (
	"cli-music-reviewer/models/entities"
	"cli-music-reviewer/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type OptionHideSplashModel struct {
	isSelected bool
	title      string
	entry      *entities.EntryRow
}

func (m *OptionHideSplashModel) Update(msg tea.Msg) (*OptionHideSplashModel, tea.Cmd) {
	return m, nil
}

func (m *OptionHideSplashModel) View() string {
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

func (m *OptionHideSplashModel) SetSelected(isSelected bool) {
	m.isSelected = isSelected
}

func NewOptionHideSplashModel(isSelected bool) *OptionHideSplashModel {
	return &OptionHideSplashModel{
		isSelected: isSelected,
	}
}
