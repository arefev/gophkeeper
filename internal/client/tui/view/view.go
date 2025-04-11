package view

import (
	"strings"
	"unicode/utf8"

	"github.com/arefev/gophkeeper/internal/client/tui/form"
	"github.com/arefev/gophkeeper/internal/client/tui/style"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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

	str := style.ErrorStyle.Render(strings.Repeat("-", count))
	str += BreakLine().One()
	str += "❗ "
	str += style.ErrorStyle.Render(err)
	str += BreakLine().One()
	str += style.ErrorStyle.Render(strings.Repeat("-", count))
	str += BreakLine().Two()

	return str
}

func Title(t string) string {
	str := BreakLine().One()
	str += style.TitleStyle.Render("## ")
	str += style.TitleStyle.Render(t)
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
	return style.HelpStyle.Render("\nCtrl+C - завершение программы")
}

func ToStart() string {
	return style.HelpStyle.Render("\nEsc - домой")
}

func Button(label string, isFocused bool) string {
	button := style.BlurredStyle.Render("[", label, "]")
	if isFocused {
		button = style.FocusedStyle.Render("[", label, "]")
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

func UpdateFocusInFields(focusIndex int, fields []*form.Input) tea.Cmd {
	cmds := make([]tea.Cmd, len(fields))
	for i, f := range fields {
		if i == focusIndex {
			cmds[i] = f.Model().Focus()
			f.Model().PromptStyle = style.FocusedStyle
			f.Model().TextStyle = style.FocusedStyle
			continue
		}
		f.Model().Blur()
		f.Model().PromptStyle = style.NoStyle
		f.Model().TextStyle = style.NoStyle
	}

	return tea.Batch(cmds...)
}
