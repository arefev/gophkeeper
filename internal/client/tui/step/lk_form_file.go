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

type lkFormFile struct {
	err        string
	app        *app.App
	fields     []*form.Input
	focusIndex int
}

func NewLKFormFile(a *app.App) *lkFormFile {
	m := lkFormFile{app: a}
	m.createFields()
	return &m
}

func (lkff *lkFormFile) createFields() {
	const (
		fieldCount     = 2
		pathCharsLimit = 500
	)
	lkff.fields = make([]*form.Input, fieldCount)
	lkff.fields[0] = form.NewInput("name", "Имя для данных", true, false)
	lkff.fields[1] = form.NewInput("path", "Путь до файла", false, false).SetCharsLimit(pathCharsLimit)
}

func (lkff *lkFormFile) WithError(err error) *lkFormFile {
	lkff.err = "Ошибка: " + err.Error()
	return lkff
}

func (lkff *lkFormFile) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, lkff.app.Conn.CheckTokenCmd)
}

func (lkff *lkFormFile) Exec() (tea.Model, tea.Cmd) {
	cmd := lkff.Init()
	return lkff, cmd
}

func (lkff *lkFormFile) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(lkff.app).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewLKTypes(lkff.app).Exec()

		case tea.KeyCtrlC:
			return lkff, tea.Quit

		case tea.KeyTab, tea.KeyShiftTab, tea.KeyUp, tea.KeyDown, tea.KeyEnter:
			if msg.Type == tea.KeyEnter && lkff.focusIndex == len(lkff.fields) {
				return NewFileSendAction(lkff.getData(), lkff.app).Exec()
			}

			if msg.Type == tea.KeyUp || msg.Type == tea.KeyShiftTab {
				lkff.focusIndex--
			} else {
				lkff.focusIndex++
			}

			if lkff.focusIndex > len(lkff.fields) {
				lkff.focusIndex = 0
			} else if lkff.focusIndex < 0 {
				lkff.focusIndex = len(lkff.fields)
			}

			return lkff, view.UpdateFocusInFields(lkff.focusIndex, lkff.fields)

		default:
			cmd := lkff.updateInputs(msg)
			return lkff, cmd
		}
	}

	cmd := lkff.updateInputs(msg)
	return lkff, cmd
}

func (lkff *lkFormFile) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(lkff.fields))

	for i := range lkff.fields {
		input, cmd := lkff.fields[i].Model().Update(msg)
		lkff.fields[i].SetModel(&input)
		cmds[i] = cmd
	}

	return tea.Batch(cmds...)
}

func (lkff *lkFormFile) View() string {
	return view.FormWithFields(
		lkff.fields,
		"Введите путь до файла для сохранения",
		form.LabelSend,
		lkff.err,
		lkff.focusIndex == len(lkff.fields),
	)
}

func (lkff *lkFormFile) getData() *model.FileData {
	data := &model.FileData{}
	for _, f := range lkff.fields {
		code := f.Code()
		switch code {
		case "name":
			data.Name = f.Model().Value()
		case "path":
			data.Path = f.Model().Value()
		}
	}

	return data
}
