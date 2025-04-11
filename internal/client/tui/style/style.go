package style

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	TitleStyle          = lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
	ErrorStyle          = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("1"))
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	HelpStyle           = BlurredStyle
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	BaseStyle           = lipgloss.NewStyle().
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("240"))
)
