package view

import (
	"strings"
	"unicode/utf8"

	"github.com/arefev/gophkeeper/internal/client/tui"
	"github.com/arefev/gophkeeper/internal/client/tui/form"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
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

	str := tui.ErrorStyle.Render(strings.Repeat("-", count))
	str += BreakLine().One()
	str += "❗ "
	str += tui.ErrorStyle.Render(err)
	str += BreakLine().One()
	str += tui.ErrorStyle.Render(strings.Repeat("-", count))
	str += BreakLine().Two()

	return str
}

func Title(t string) string {
	str := BreakLine().One()
	str += tui.TitleStyle.Render("## ")
	str += tui.TitleStyle.Render(t)
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
	return tui.HelpStyle.Render("\nCtrl+C - завершение программы")
}

func ToStart() string {
	return tui.HelpStyle.Render("\nEsc - домой")
}

func Button(label string, isFocused bool) string {
	button := tui.BlurredStyle.Render("[", label, "]")
	if isFocused {
		button = tui.FocusedStyle.Render("[", label, "]")
	}

	str := BreakLine().Two()
	str += button
	str += BreakLine().One()

	return str
}

func FormWithFields(fields []*form.Input, title, btnLabel, err string, isBtnFocused bool) string {
	str := Title(title)
	if err != "" {
		str += Error(err)
	}

	for i := range fields {
		str += fields[i].Model().View()
		if i < len(fields)-1 {
			str += BreakLine().One()
		}
	}

	str += Button(btnLabel, isBtnFocused)
	str += Quit() + ToStart()
	return str
}
