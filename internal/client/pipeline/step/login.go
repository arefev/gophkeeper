package step

import (
	"github.com/arefev/gophkeeper/internal/client/pipeline/view"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type login struct {
	focusIndex int
	inputs     []textinput.Model
	err        string
}

func NewLogin() *login {
	m := login{
		inputs: make([]textinput.Model, 2),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = view.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Логин"
			t.Focus()
			t.PromptStyle = view.FocusedStyle
			t.TextStyle = view.FocusedStyle
		case 1:
			t.Placeholder = "Пароль"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return &m
}

func (l *login) WithError(err error) *login {
	l.err = "Ошибка: " + err.Error()
	return l
}

func (l *login) Init() tea.Cmd {
	return textinput.Blink
}

func (l *login) Exec() (tea.Model, tea.Cmd) {
	return l, l.Init()
}

func (l *login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			scr := NewStart()
			return scr, scr.Init()
		case tea.KeyCtrlC:
			return l, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			if msg.Type == tea.KeyEnter && l.focusIndex == len(l.inputs) {
				return NewLoginAction().Exec()
			}

			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
				l.focusIndex--
			} else {
				l.focusIndex++
			}

			if l.focusIndex > len(l.inputs) {
				l.focusIndex = 0
			} else if l.focusIndex < 0 {
				l.focusIndex = len(l.inputs)
			}

			cmds := make([]tea.Cmd, len(l.inputs))
			for i := 0; i <= len(l.inputs)-1; i++ {
				if i == l.focusIndex {
					cmds[i] = l.inputs[i].Focus()
					l.inputs[i].PromptStyle = view.FocusedStyle
					l.inputs[i].TextStyle = view.FocusedStyle
					continue
				}
				l.inputs[i].Blur()
				l.inputs[i].PromptStyle = view.NoStyle
				l.inputs[i].TextStyle = view.NoStyle
			}

			return l, tea.Batch(cmds...)
		}
	}

	cmd := l.updateInputs(msg)
	return l, cmd
}

func (l *login) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(l.inputs))

	for i := range l.inputs {
		l.inputs[i], cmds[i] = l.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (l *login) View() string {
	str := view.Title("Авторизация")
	if l.err != "" {
		str += view.Error(l.err)
	}

	for i := range l.inputs {
		str += l.inputs[i].View()
		if i < len(l.inputs)-1 {
			str += view.Break(1)
		}
	}

	str += view.Button("Войти", l.focusIndex == len(l.inputs))
	str += view.Quit() + view.ToStart()
	return str
}
