package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedColor = "205"
	blurredColor = "240"
	titleColor   = "111"
	errorColor   = "1"
	cursorColor  = "244"
)

var (
	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(focusedColor))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color(blurredColor))
	TitleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color(titleColor))
	ErrorStyle          = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color(errorColor))
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(cursorColor))
	BorderStyle         = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color(blurredColor))
)
