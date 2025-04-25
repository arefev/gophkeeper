package step

import (
	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/form"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type lkFormCreds struct {
	err        string
	app        *app.App
	fields     []*form.Input
	focusIndex int
}

func NewLKFormCreds(a *app.App) *lkFormCreds {
	m := lkFormCreds{app: a}
	m.createFields()
	return &m
}

func (lkfc *lkFormCreds) createFields() {
	const fieldCount = 3
	lkfc.fields = make([]*form.Input, fieldCount)
	lkfc.fields[0] = form.NewInput(form.CodeName, "Имя для данных", true, false)
	lkfc.fields[1] = form.NewInput(form.CodeLogin, "Логин", false, false)
	lkfc.fields[2] = form.NewInput(form.CodePwd, "Пароль", false, true)
}

func (lkfc *lkFormCreds) WithError(err error) *lkFormCreds {
	lkfc.err = "Ошибка: " + err.Error()
	return lkfc
}

func (lkfc *lkFormCreds) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, lkfc.app.Conn.CheckTokenCmd)
}

func (lkfc *lkFormCreds) Exec() (tea.Model, tea.Cmd) {
	cmd := lkfc.Init()
	return lkfc, cmd
}

func (lkfc *lkFormCreds) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(lkfc.app).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewLKTypes(lkfc.app).Exec()

		case tea.KeyCtrlC:
			return lkfc, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			if msg.Type == tea.KeyEnter && lkfc.focusIndex == len(lkfc.fields) {
				return NewCredsSendAction(lkfc.getData(), lkfc.app).Exec()
			}

			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
				lkfc.focusIndex--
			} else {
				lkfc.focusIndex++
			}

			if lkfc.focusIndex > len(lkfc.fields) {
				lkfc.focusIndex = 0
			} else if lkfc.focusIndex < 0 {
				lkfc.focusIndex = len(lkfc.fields)
			}

			return lkfc, view.UpdateFocusInFields(lkfc.focusIndex, lkfc.fields)

		default:
			cmd := lkfc.updateInputs(msg)
			return lkfc, cmd
		}
	}

	cmd := lkfc.updateInputs(msg)
	return lkfc, cmd
}

func (lkfc *lkFormCreds) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(lkfc.fields))

	for i := range lkfc.fields {
		input, cmd := lkfc.fields[i].Model().Update(msg)
		lkfc.fields[i].SetModel(&input)
		cmds[i] = cmd
	}

	return tea.Batch(cmds...)
}

func (lkfc *lkFormCreds) View() string {
	return view.FormWithFields(
		lkfc.fields,
		"Введите логин/пароль для сохранения",
		"Отправить",
		lkfc.err,
		lkfc.focusIndex == len(lkfc.fields),
	)
}

func (lkfc *lkFormCreds) getData() *model.CredsData {
	data := &model.CredsData{}
	for _, f := range lkfc.fields {
		code := f.Code()
		switch code {
		case form.CodeLogin:
			data.Login = f.Model().Value()
		case form.CodePwd:
			data.Password = f.Model().Value()
		case form.CodeName:
			data.Name = f.Model().Value()
		}
	}

	return data
}
