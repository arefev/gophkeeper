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

type BankSendActionSuccess struct {
}

type BankSendActionFail struct {
	Err error
}

type bankSendAction struct {
	app     *app.App
	data    *model.BankData
	spinner spinner.Model
}

func NewBankSendAction(data *model.BankData, a *app.App) *bankSendAction {
	return &bankSendAction{
		spinner: view.Spinner(),
		data:    data,
		app:     a,
	}
}

func (bsa *bankSendAction) ActionCmd() tea.Msg {
	ctx := context.Background()
	data := fmt.Sprintf("Number: %s\nExpired: %s\nCVV: %s", bsa.data.Number, bsa.data.Exp, bsa.data.CVV)
	err := bsa.app.Conn.TextUpload(ctx, []byte(data), bsa.data.Name, "card")
	if err != nil {
		bsa.app.Log.Error("TextUpload failed", zap.Error(err))
		return BankSendActionFail{Err: errors.New("не удалось сохранить данные")}
	}

	return BankSendActionSuccess{}
}

func (bsa *bankSendAction) Init() tea.Cmd {
	return tea.Batch(bsa.spinner.Tick, bsa.ActionCmd, bsa.app.Conn.CheckTokenCmd)
}

func (bsa *bankSendAction) Exec() (tea.Model, tea.Cmd) {
	cmd := bsa.Init()
	return bsa, cmd
}

func (bsa *bankSendAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(bsa.app).Exec()

	case BankSendActionSuccess:
		return NewLKTypes(bsa.app).WithSuccess().Exec()

	case BankSendActionFail:
		return NewLogin(bsa.app).WithError(msg.Err).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return bsa, tea.Quit

		default:
			return bsa, nil
		}

	default:
		var cmd tea.Cmd
		bsa.spinner, cmd = bsa.spinner.Update(msg)
		return bsa, cmd
	}
}

func (bsa *bankSendAction) View() string {
	str := view.BreakLine().Two()
	str += bsa.spinner.View()
	str += " Минутку..." + view.BreakLine().One()
	str += view.Quit()
	return str
}
