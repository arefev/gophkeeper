package step

import (
	"errors"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type DeleteActionSuccess struct {
}

type DeleteActionFail struct {
	Err error
}

type deleteAction struct {
	uuid    string
	app     *app.App
	spinner spinner.Model
}

func NewDeleteAction(uuid string, a *app.App) *deleteAction {
	return &deleteAction{
		spinner: view.Spinner(),
		uuid:    uuid,
		app:     a,
	}
}

func (del *deleteAction) ActionCmd() tea.Msg {
	return DeleteActionFail{Err: errors.New("не получилось удалить файл")}
}

func (del *deleteAction) Init() tea.Cmd {
	return tea.Batch(del.spinner.Tick, del.ActionCmd, del.app.Conn.CheckTokenCmd)
}

func (del *deleteAction) Exec() (tea.Model, tea.Cmd) {
	cmd := del.Init()
	return del, cmd
}

func (del *deleteAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(del.app).Exec()

	case DeleteActionSuccess:
		return NewLKTypes(del.app).WithSuccess().Exec()

	case DeleteActionFail:
		return NewLKList(del.app).WithError(msg.Err).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return del, tea.Quit

		default:
			return del, nil
		}

	default:
		var cmd tea.Cmd
		del.spinner, cmd = del.spinner.Update(msg)
		return del, cmd
	}
}

func (del *deleteAction) View() string {
	str := view.BreakLine().Two()
	str += del.spinner.View()
	str += " Минутку..." + view.BreakLine().One()
	str += view.Quit()
	return str
}
