package step

import (
	"github.com/arefev/gophkeeper/internal/client/tui"
	"github.com/arefev/gophkeeper/internal/client/tui/form"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type login struct {
	err        string
	fields     []*form.Input
	focusIndex int
}

func NewLogin() *login {
	m := login{}
	m.createFields()
	return &m
}

func (l *login) createFields() {
	const fieldCount = 2
	l.fields = make([]*form.Input, fieldCount)
	l.fields[0] = form.NewInput("login", "Логин", true, false)
	l.fields[1] = form.NewInput("pwd", "Пароль", false, true)
}

func (l *login) WithError(err error) *login {
	l.err = "Ошибка: " + err.Error()
	return l
}

func (l *login) Init() tea.Cmd {
	return textinput.Blink
}

func (l *login) Exec() (tea.Model, tea.Cmd) {
	cmd := l.Init()
	return l, cmd
}

func (l *login) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyEsc:
			return NewStart().Exec()
		case tea.KeyCtrlC:
			return l, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			if msg.Type == tea.KeyEnter && l.focusIndex == len(l.fields) {
				return NewLoginAction(l.getAuthData()).Exec()
			}

			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
				l.focusIndex--
			} else {
				l.focusIndex++
			}

			if l.focusIndex > len(l.fields) {
				l.focusIndex = 0
			} else if l.focusIndex < 0 {
				l.focusIndex = len(l.fields)
			}

			cmds := make([]tea.Cmd, len(l.fields))
			for i := 0; i <= len(l.fields)-1; i++ {
				if i == l.focusIndex {
					cmds[i] = l.fields[i].Model().Focus()
					l.fields[i].Model().PromptStyle = tui.FocusedStyle
					l.fields[i].Model().TextStyle = tui.FocusedStyle
					continue
				}
				l.fields[i].Model().Blur()
				l.fields[i].Model().PromptStyle = tui.NoStyle
				l.fields[i].Model().TextStyle = tui.NoStyle
			}

			return l, tea.Batch(cmds...)
		}
	}

	cmd := l.updateInputs(msg)
	return l, cmd
}

func (l *login) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(l.fields))

	for i := range l.fields {
		input, cmd := l.fields[i].Model().Update(msg)
		l.fields[i].SetModel(&input)
		cmds[i] = cmd
	}

	return tea.Batch(cmds...)
}

func (l *login) View() string {
	return view.FormWithFields(l.fields, "Авторизация", "Войти", l.err, l.focusIndex == len(l.fields))
}

func (l *login) getAuthData() *model.LoginData {
	data := &model.LoginData{}
	for _, f := range l.fields {
		code := f.Code()
		switch code {
		case "login":
			data.Login = f.Model().Value()
		case "pwd":
			data.Password = f.Model().Value()
		}
	}

	return data
}
