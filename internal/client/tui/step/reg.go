package step

import (
	"github.com/arefev/gophkeeper/internal/client/tui"
	"github.com/arefev/gophkeeper/internal/client/tui/form"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type reg struct {
	err        string
	fields     []*form.Input
	focusIndex int
}

func NewReg() *reg {
	m := reg{}
	m.createFields()
	return &m
}

func (r *reg) createFields() {
	const fieldCount = 3
	r.fields = make([]*form.Input, fieldCount)
	r.fields[0] = form.NewInput("login", "Логин", true, false)
	r.fields[1] = form.NewInput("pwd", "Пароль", false, true)
	r.fields[2] = form.NewInput("pwdConfirm", "Повторите пароль", false, true)
}

func (r *reg) WithError(err error) *reg {
	r.err = "Ошибка: " + err.Error()
	return r
}

func (r *reg) Init() tea.Cmd {
	return textinput.Blink
}

func (r *reg) Exec() (tea.Model, tea.Cmd) {
	return r, r.Init()
}

func (r *reg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyEsc:
			return NewStart().Exec()
		case tea.KeyCtrlC:
			return r, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			if msg.Type == tea.KeyEnter && r.focusIndex == len(r.fields) {
				return NewRegAction().Exec()
			}

			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
				r.focusIndex--
			} else {
				r.focusIndex++
			}

			if r.focusIndex > len(r.fields) {
				r.focusIndex = 0
			} else if r.focusIndex < 0 {
				r.focusIndex = len(r.fields)
			}

			cmds := make([]tea.Cmd, len(r.fields))
			for i := 0; i <= len(r.fields)-1; i++ {
				if i == r.focusIndex {
					cmds[i] = r.fields[i].Model().Focus()
					r.fields[i].Model().PromptStyle = tui.FocusedStyle
					r.fields[i].Model().TextStyle = tui.FocusedStyle
					continue
				}
				r.fields[i].Model().Blur()
				r.fields[i].Model().PromptStyle = tui.NoStyle
				r.fields[i].Model().TextStyle = tui.NoStyle
			}

			return r, tea.Batch(cmds...)
		}
	}

	cmd := r.updateInputs(msg)
	return r, cmd
}

func (r *reg) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(r.fields))

	for i := range r.fields {
		input, cmd := r.fields[i].Model().Update(msg)
		r.fields[i].SetModel(&input)
		cmds[i] = cmd
	}

	return tea.Batch(cmds...)
}

func (r *reg) View() string {
	return view.FormWithFields(r.fields, "Регистрация", "Отправить", r.err, r.focusIndex == len(r.fields))
}
