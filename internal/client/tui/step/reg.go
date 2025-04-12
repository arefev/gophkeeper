package step

import (
	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/tui/form"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type reg struct {
	err        string
	app        *app.App
	fields     []*form.Input
	focusIndex int
}

func NewReg(a *app.App) *reg {
	m := reg{app: a}
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
	cmd := r.Init()
	return r, cmd
}

func (r *reg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyEsc:
			return NewStart(r.app).Exec()

		case tea.KeyCtrlC:
			return r, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			if msg.Type == tea.KeyEnter && r.focusIndex == len(r.fields) {
				return NewRegAction(r.app).Exec()
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

			return r, view.UpdateFocusInFields(r.focusIndex, r.fields)

		default:
			cmd := r.updateInputs(msg)
			return r, cmd
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
