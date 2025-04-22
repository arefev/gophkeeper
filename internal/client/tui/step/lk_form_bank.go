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

type lkFormBank struct {
	err        string
	app        *app.App
	fields     []*form.Input
	focusIndex int
}

func NewLKFormBank(a *app.App) *lkFormBank {
	m := lkFormBank{app: a}
	m.createFields()
	return &m
}

func (lkfb *lkFormBank) createFields() {
	const fieldCount = 4
	lkfb.fields = make([]*form.Input, fieldCount)
	lkfb.fields[0] = form.NewInput("name", "Имя для данных", true, false)
	lkfb.fields[1] = form.NewInput("number", "Номер карты", false, false)
	lkfb.fields[2] = form.NewInput("exp", "Срок действия", false, false)
	lkfb.fields[3] = form.NewInput("cvv", "CVV", false, false)
}

func (lkfb *lkFormBank) WithError(err error) *lkFormBank {
	lkfb.err = "Ошибка: " + err.Error()
	return lkfb
}

func (lkfb *lkFormBank) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, lkfb.app.Conn.CheckTokenCmd)
}

func (lkfb *lkFormBank) Exec() (tea.Model, tea.Cmd) {
	cmd := lkfb.Init()
	return lkfb, cmd
}

func (lkfb *lkFormBank) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(lkfb.app).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewLKTypes(lkfb.app).Exec()

		case tea.KeyCtrlC:
			return lkfb, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			if msg.Type == tea.KeyEnter && lkfb.focusIndex == len(lkfb.fields) {
				return NewBankSendAction(lkfb.getData(), lkfb.app).Exec()
			}

			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
				lkfb.focusIndex--
			} else {
				lkfb.focusIndex++
			}

			if lkfb.focusIndex > len(lkfb.fields) {
				lkfb.focusIndex = 0
			} else if lkfb.focusIndex < 0 {
				lkfb.focusIndex = len(lkfb.fields)
			}

			return lkfb, view.UpdateFocusInFields(lkfb.focusIndex, lkfb.fields)

		default:
			cmd := lkfb.updateInputs(msg)
			return lkfb, cmd
		}
	}

	cmd := lkfb.updateInputs(msg)
	return lkfb, cmd
}

func (lkfb *lkFormBank) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(lkfb.fields))

	for i := range lkfb.fields {
		input, cmd := lkfb.fields[i].Model().Update(msg)
		lkfb.fields[i].SetModel(&input)
		cmds[i] = cmd
	}

	return tea.Batch(cmds...)
}

func (lkfb *lkFormBank) View() string {
	return view.FormWithFields(
		lkfb.fields,
		"Введите данные банковской карты для сохранения",
		"Отправить",
		lkfb.err,
		lkfb.focusIndex == len(lkfb.fields),
	)
}

func (lkfb *lkFormBank) getData() *model.BankData {
	data := &model.BankData{}
	for _, f := range lkfb.fields {
		code := f.Code()
		switch code {
		case "number":
			data.Number = f.Model().Value()
		case "exp":
			data.Exp = f.Model().Value()
		case "cvv":
			data.CVV = f.Model().Value()
		case "name":
			data.Name = f.Model().Value()
		}
	}

	return data
}
