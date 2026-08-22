package components

import (
	"cli-music-reviewer/config"
	"cli-music-reviewer/events"
	"cli-music-reviewer/interfaces"
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/styles"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EntryBrowserModel struct {
	activeIndex  int
	children     []*EntryRowModel
	showControls bool
}

func (m *EntryBrowserModel) Update(msg tea.Msg) (*EntryBrowserModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, func() tea.Msg { return events.EntryCreateRequestedMsg{} }
		case "enter":
			return m, func() tea.Msg { return events.EntryEditRequestedMsg{Index: m.activeIndex} }
		}
	}

	for i := range m.children {
		m.children[i].SetSelected(i == m.activeIndex)
		m.children[i], cmd = m.children[i].Update(msg)
	}
	return m, cmd
}

func (m *EntryBrowserModel) View() string {
	header := styles.ConfigHeaderStyle.Render(fmt.Sprintf(" %s ", config.AppTitle))
	versionPart := lipgloss.NewStyle().Foreground(styles.Gray).Render(" Configuration ")
	headerRow := lipgloss.JoinHorizontal(lipgloss.Bottom, header, versionPart)

	var rows []string
	rows = append(rows, headerRow, "")

	for _, entry := range m.children {
		rows = append(rows, entry.View())
	}

	content := strings.Join(rows, "\n")
	return styles.ConfigHeroStyle.Render(content)
}

func (m *EntryBrowserModel) CursorUp() {
	if m.activeIndex <= 0 {
		m.activeIndex = 0
		return
	}
	m.activeIndex--
}

func (m *EntryBrowserModel) CursorDown() {
	if m.activeIndex >= len(m.children)-1 {
		m.activeIndex = len(m.children) - 1
		return
	}
	m.activeIndex++
}

func NewEntryBrowser(showControls bool, repos *repositories.AppRepositories) *EntryBrowserModel {
	entryRows := repos.EntryRowRepository.GetActiveRows()

	var rowModels []*EntryRowModel
	for _, row := range entryRows {
		rowModels = append(rowModels, &EntryRowModel{
			isSelected: false,
			title:      row.Title,
			timestamp:  row.UpdatedAt,
		})
	}

	return &EntryBrowserModel{
		activeIndex:  0,
		children:     rowModels,
		showControls: showControls,
	}
}

var (
	_ interfaces.ComponentInterface = (*EntryBrowserModel)(nil)
)
