package step

import (
	"context"
	"errors"
	"fmt"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type CredsSendActionSuccess struct {
}

type CredsSendActionFail struct {
	Err error
}

type credsSendAction struct {
	app     *app.App
	data    *model.CredsData
	spinner spinner.Model
}

func NewCredsSendAction(data *model.CredsData, a *app.App) *credsSendAction {
	return &credsSendAction{
		spinner: view.Spinner(),
		data:    data,
		app:     a,
	}
}

func (csa *credsSendAction) ActionCmd() tea.Msg {
	ctx := context.Background()
	// TODO: validation needed
	data := fmt.Sprintf("Login: %s\nPassword: %s", csa.data.Login, csa.data.Password)
	err := csa.app.Conn.TextUpload(ctx, []byte(data), csa.data.Name, "creds")
	if err != nil {
		csa.app.Log.Error("TextUpload failed", zap.Error(err))
		return CredsSendActionFail{Err: errors.New("не удалось сохранить данные")}
	}

	return CredsSendActionSuccess{}
}

func (csa *credsSendAction) Init() tea.Cmd {
	return tea.Batch(csa.spinner.Tick, csa.ActionCmd, csa.app.Conn.CheckTokenCmd)
}

func (csa *credsSendAction) Exec() (tea.Model, tea.Cmd) {
	cmd := csa.Init()
	return csa, cmd
}

func (csa *credsSendAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(csa.app).Exec()

	case CredsSendActionSuccess:
		return NewLKTypes(csa.app).WithSuccess().Exec()

	case CredsSendActionFail:
		return NewLogin(csa.app).WithError(msg.Err).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return csa, tea.Quit

		default:
			return csa, nil
		}

	default:
		var cmd tea.Cmd
		csa.spinner, cmd = csa.spinner.Update(msg)
		return csa, cmd
	}
}

func (csa *credsSendAction) View() string {
	str := view.BreakLine().Two()
	str += csa.spinner.View()
	str += " Минутку..." + view.BreakLine().One()
	str += view.Quit()
	return str
}
