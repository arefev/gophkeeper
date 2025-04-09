package step

import (
	"fmt"

	"github.com/arefev/gophkeeper/internal/client/pipeline/view"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type reg struct {
	focusIndex int
	inputs     []textinput.Model
	err        string
}

func NewReg() reg {
	m := reg{
		inputs: make([]textinput.Model, 3),
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
		case 2:
			t.Placeholder = "Еще раз пароль"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (r reg) WithError(err string) reg {
	r.err = fmt.Sprintf("Ошибка: %s", err)
	return r
}

func (r reg) Init() tea.Cmd {
	return textinput.Blink
}

func (r reg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			scr := NewStart()
			return scr, scr.Init()
		case tea.KeyCtrlC:
			return r, tea.Quit
		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			if msg.Type == tea.KeyEnter && r.focusIndex == len(r.inputs) {
				pr := NewRegProcessing()
				return pr, pr.Init()
			}

			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
				r.focusIndex--
			} else {
				r.focusIndex++
			}

			if r.focusIndex > len(r.inputs) {
				r.focusIndex = 0
			} else if r.focusIndex < 0 {
				r.focusIndex = len(r.inputs)
			}

			cmds := make([]tea.Cmd, len(r.inputs))
			for i := 0; i <= len(r.inputs)-1; i++ {
				if i == r.focusIndex {
					cmds[i] = r.inputs[i].Focus()
					r.inputs[i].PromptStyle = view.FocusedStyle
					r.inputs[i].TextStyle = view.FocusedStyle
					continue
				}
				r.inputs[i].Blur()
				r.inputs[i].PromptStyle = view.NoStyle
				r.inputs[i].TextStyle = view.NoStyle
			}

			return r, tea.Batch(cmds...)
		}
	}

	cmd := r.updateInputs(msg)
	return r, cmd
}

func (r *reg) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(r.inputs))

	for i := range r.inputs {
		r.inputs[i], cmds[i] = r.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (r reg) View() string {
	str := view.Title("Регистрация")
	if r.err != "" {
		str += view.Error(r.err)
	}

	for i := range r.inputs {
		str += r.inputs[i].View()
		if i < len(r.inputs)-1 {
			str += view.Break(1)
		}
	}

	str += view.Button("Отправить", r.focusIndex == len(r.inputs))
	str += view.Quit() + view.ToStart()
	return str
}
