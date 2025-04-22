package step

import (
	"context"
	"errors"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type FileSendActionSuccess struct {
}

type FileSendActionFail struct {
	Err error
}

type fileSendAction struct {
	app     *app.App
	data    *model.FileData
	spinner spinner.Model
}

func NewFileSendAction(data *model.FileData, a *app.App) *fileSendAction {
	return &fileSendAction{
		spinner: view.Spinner(),
		data:    data,
		app:     a,
	}
}

func (fsa *fileSendAction) ActionCmd() tea.Msg {
	ctx := context.Background()
	// TODO: validation needed
	err := fsa.app.Conn.FileUpload(ctx, fsa.data.Path, fsa.data.Name, "file")
	if err != nil {
		fsa.app.Log.Error("FileUpload failed", zap.Error(err))
		return FileSendActionFail{Err: errors.New("не удалось сохранить данные")}
	}

	return FileSendActionSuccess{}
}

func (fsa *fileSendAction) Init() tea.Cmd {
	return tea.Batch(fsa.spinner.Tick, fsa.ActionCmd, fsa.app.Conn.CheckTokenCmd)
}

func (fsa *fileSendAction) Exec() (tea.Model, tea.Cmd) {
	cmd := fsa.Init()
	return fsa, cmd
}

func (fsa *fileSendAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(fsa.app).Exec()

	case FileSendActionSuccess:
		return NewLKTypes(fsa.app).WithSuccess().Exec()

	case FileSendActionFail:
		return NewLogin(fsa.app).WithError(msg.Err).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return fsa, tea.Quit

		default:
			return fsa, nil
		}

	default:
		var cmd tea.Cmd
		fsa.spinner, cmd = fsa.spinner.Update(msg)
		return fsa, cmd
	}
}

func (fsa *fileSendAction) View() string {
	str := view.BreakLine().Two()
	str += fsa.spinner.View()
	str += " Минутку..." + view.BreakLine().One()
	str += view.Quit()
	return str
}
