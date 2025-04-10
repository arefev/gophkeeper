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

type breakLine struct{}

func BreakLine() *breakLine {
	return &breakLine{}
}

func (bl *breakLine) One() string {
	return "\n"
}

func (bl *breakLine) Two() string {
	return "\n\n"
}

func Error(err string) string {
	const emLen = 3
	count := utf8.RuneCountInString(err) + emLen

	str := ErrorStyle.Render(strings.Repeat("-", count))
	str += BreakLine().One()
	str += "❗ "
	str += ErrorStyle.Render(err)
	str += BreakLine().One()
	str += ErrorStyle.Render(strings.Repeat("-", count))
	str += BreakLine().Two()

	return str
}

func Title(t string) string {
	str := BreakLine().One()
	str += TitleStyle.Render("## ")
	str += TitleStyle.Render(t)
	str += BreakLine().Two()

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

	str := BreakLine().Two()
	str += button
	str += BreakLine().One()

	return str
}
