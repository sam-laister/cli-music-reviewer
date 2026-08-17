package components

import (
	"cli-music-reviewer/config"
	"cli-music-reviewer/styles"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EntryBrowserModel struct {
	activeIndex int
	children    []*EntryRowModel
}

func (m *EntryBrowserModel) Update(msg tea.Msg) (*EntryBrowserModel, tea.Cmd) {
	var cmd tea.Cmd

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

func NewEntryBrowser() *EntryBrowserModel {
	return &EntryBrowserModel{
		activeIndex: 0,
		children: []*EntryRowModel{
			NewEntryRow(true, "Test Row", "17/08/26"),
			NewEntryRow(false, "Test Row 2", "17/08/26"),
			NewEntryRow(false, "Test Row 3", "17/08/26"),
		},
	}
}
