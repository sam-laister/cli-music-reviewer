package components

import (
	"cli-music-reviewer/models/entities"
	"cli-music-reviewer/styles"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type EntryRowModel struct {
	isSelected bool
	id         uint64
	title      string
	timestamp  string
	entry      *entities.EntryRow
}

func (m *EntryRowModel) Update(msg tea.Msg) (*EntryRowModel, tea.Cmd) {
	return m, nil
}

func (m *EntryRowModel) View() string {
	cursor := "  "
	label := m.title

	value := fmt.Sprintf("[ %s ]", m.timestamp)

	if m.isSelected {
		cursor = styles.CursorStyle.Render("> ")
		label = styles.SelectedItemStyle.Render(label)
		value = styles.SelectedItemStyle.Render(value)
	}

	return fmt.Sprintf("%s%s %s", cursor, label, value)
}

func (m *EntryRowModel) SetSelected(isSelected bool) {
	m.isSelected = isSelected
}

func NewEntryRow(isSelected bool, entry *entities.EntryRow, timestamp string) *EntryRowModel {
	return &EntryRowModel{
		isSelected: isSelected,
		id:         entry.ID,
		title:      entry.Title,
		timestamp:  timestamp,
		entry:      entry,
	}
}
