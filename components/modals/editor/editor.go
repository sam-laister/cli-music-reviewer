package editor

import (
	"cli-music-reviewer/events"
	"cli-music-reviewer/models/entities"
	"cli-music-reviewer/repositories"
	"cli-music-reviewer/styles"

	"charm.land/glamour/v2"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	defaultEditorWidth  = 100
	defaultEditorHeight = 24

	// paneWidthOverhead: two panes side by side, each with a 1-char border
	// on both sides (lipgloss adds border on top of the requested Width) —
	// 2 panes * 2 border chars = 4.
	paneWidthOverhead = 4

	// paneHeightOverhead accounts for every line surrounding the pane block
	// that isn't part of it, so the whole frame fits within the terminal's
	// reported height without Bubble Tea silently chopping lines off the
	// top (it can't scroll into history in alt-screen mode): pane border
	// (top+bottom, +2) + this view's own header/blank/blank/instructions
	// lines (+4). The editor is full-bleed — homepage.go doesn't wrap it in
	// anything else, unlike the centered "central modal" views.
	paneHeightOverhead = 6
)

type Model struct {
	entry        *entities.EntryRow
	originalBody string

	editor   textarea.Model
	preview  viewport.Model
	renderer *glamour.TermRenderer

	confirmingDiscard bool

	width  int
	height int
}

func New(entryID uint64, repo repositories.EntryRowRepositoryInterface, width, height int) (*Model, error) {
	entry, err := repo.FindByID(entryID)
	if err != nil {
		return nil, err
	}

	if width <= 0 {
		width = defaultEditorWidth
	}
	if height <= 0 {
		height = defaultEditorHeight
	}

	editor := textarea.New()
	editor.Placeholder = "Write your review in Markdown…"
	editor.SetValue(entry.Body)
	editor.Focus()

	m := &Model{
		entry:        entry,
		originalBody: entry.Body,
		editor:       editor,
		preview:      viewport.New(0, 0),
	}
	m.resize(width, height)

	return m, nil
}

func (m *Model) paneWidth() int {
	return (m.width - paneWidthOverhead) / 2
}

func (m *Model) paneHeight() int {
	return m.height - paneHeightOverhead
}

func (m *Model) resize(width, height int) {
	m.width = width
	m.height = height

	paneWidth := m.paneWidth()
	paneHeight := m.paneHeight()

	m.editor.SetWidth(paneWidth)
	m.editor.SetHeight(paneHeight)
	m.preview.Width = paneWidth
	m.preview.Height = paneHeight

	renderer, err := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(paneWidth),
	)
	if err == nil {
		m.renderer = renderer
	}

	m.updatePreview()
}

func (m *Model) updatePreview() {
	if m.renderer == nil {
		m.preview.SetContent(m.editor.Value())
		return
	}

	rendered, err := m.renderer.Render(m.editor.Value())
	if err != nil {
		m.preview.SetContent(m.editor.Value())
		return
	}
	m.preview.SetContent(rendered)
}

func (m *Model) save() tea.Msg {
	m.entry.Body = m.editor.Value()
	m.entry.MarkUpdated()
	return events.ReviewSaveRequestedMsg{Entry: m.entry}
}

func (m *Model) dirty() bool {
	return m.editor.Value() != m.originalBody
}

func (m *Model) Update(msg tea.Msg) (*Model, tea.Cmd) {
	if m.confirmingDiscard {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			if keyMsg.String() == "y" {
				return m, func() tea.Msg { return events.ReviewEditCancelledMsg{} }
			}
			m.confirmingDiscard = false
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+x":
			if m.dirty() {
				m.confirmingDiscard = true
				return m, nil
			}
			return m, func() tea.Msg { return events.ReviewEditCancelledMsg{} }
		case "ctrl+s":
			return m, m.save
		}
	}

	var cmd tea.Cmd
	m.editor, cmd = m.editor.Update(msg)
	m.updatePreview()
	return m, cmd
}

func (m *Model) View() string {
	paneStyle := lipgloss.NewStyle().
		Width(m.paneWidth()).
		Height(m.paneHeight()).
		Padding(0, 1).
		BorderStyle(lipgloss.NormalBorder())

	editorPane := paneStyle.BorderForeground(styles.Magenta).Render(m.editor.View())
	previewPane := paneStyle.BorderForeground(styles.Cyan).Render(m.preview.View())

	header := styles.ConfigHeaderStyle.Render(" Review: " + m.entry.Title + " ")
	body := lipgloss.JoinHorizontal(lipgloss.Top, editorPane, previewPane)

	instructions := styles.InstructionStyle.Render("ctrl+s to save • ctrl+x to close")
	if m.confirmingDiscard {
		instructions = styles.ErrorStyle.Render("Discard unsaved changes? y to confirm • any other key to keep editing")
	}

	return header + "\n\n" + body + "\n\n" + instructions
}
