package styles

import "github.com/charmbracelet/lipgloss"

var (
	Magenta = lipgloss.Color("#fe0979")
	Cyan    = lipgloss.Color("#00f2fe")
	White   = lipgloss.Color("#ffffff")
	Gray    = lipgloss.Color("240")

	AsciiHeroStyle = lipgloss.NewStyle().
			Margin(1, 0, 2, 0).
			Padding(1, 4).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(Cyan).
			Align(lipgloss.Center)

	AsciiTitleStyle = lipgloss.NewStyle().Foreground(Magenta).Bold(true)
	SubtitleStyle   = lipgloss.NewStyle().Foreground(Cyan).Italic(true).MarginTop(1)

	ConfigHeroStyle = lipgloss.NewStyle().
			Margin(1, 0, 1, 0).
			Padding(1, 3).
			Background(lipgloss.Color("236")).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Magenta)

	ConfigHeaderStyle = lipgloss.NewStyle().Foreground(White).Background(Magenta).Padding(0, 1).Bold(true)
	CursorStyle       = lipgloss.NewStyle().Foreground(Magenta).Bold(true)
	SelectedItemStyle = lipgloss.NewStyle().Foreground(Cyan).Bold(true)
	InstructionStyle  = lipgloss.NewStyle().Foreground(Gray)

	ModalStyle = lipgloss.NewStyle().
			Margin(1, 0, 1, 0).
			Padding(1, 3).
			Background(lipgloss.Color("236")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(Cyan)
)
