package style

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	FocusedColor            = "205"
	BlurredColor            = "240"
	TitleColor              = "111"
	ErrorColor              = "1"
	SuccessColor            = "2"
	CursorColor             = "244"
	SelectedForegroundColor = "229"
	SelectedBackgroundColor = "57"
	ColumnWidthS            = 5
	ColumnWidthM            = 10
	ColumnWitdthL           = 20
	TableHeight             = 7
)

var (
	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(FocusedColor))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(BlurredColor))
	TitleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color(TitleColor))
	ErrorStyle          = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(ErrorColor))
	SuccessStyle        = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(SuccessColor))
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(CursorColor))
	BorderStyle         = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(BlurredColor))
)
