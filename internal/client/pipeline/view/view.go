package view

import (
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/spinner"
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
)

func Break(count int) string {
	if count <= 0 {
		count = 1
	}
	return strings.Repeat("\n", count)
}

func Error(err string) string {
	count := utf8.RuneCountInString(err) + 3

	str := ErrorStyle.Render(strings.Repeat("-", count))
	str += Break(1)
	str += "❗ "
	str += ErrorStyle.Render(err)
	str += Break(1)
	str += ErrorStyle.Render(strings.Repeat("-", count))
	str += Break(2)

	return str
}

func Title(t string) string {
	str := Break(1)
	str += TitleStyle.Render("## ")
	str += TitleStyle.Render(t)
	str += Break(2)

	return str
}

func Spinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Pulse
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return s
}

func Quit() string {
	return HelpStyle.Render("\nCtrl+C - завершение программы")
}

func ToStart() string {
	return HelpStyle.Render("\nEsc - домой")
}

func Button(label string, isFocused bool) string {
	button := BlurredStyle.Render("[", label, "]")
	if isFocused {
		button = FocusedStyle.Render("[", label, "]")
	}

	str := Break(2)
	str += button
	str += Break(1)

	return str
}
