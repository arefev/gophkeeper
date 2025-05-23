package step

import (
	"context"
	"errors"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type RegActionSuccess struct {
}

type RegActionFail struct {
	Err error
}

type regAction struct {
	app     *app.App
	regData *model.RegData
	spinner spinner.Model
}

func NewRegAction(a *app.App, data *model.RegData) *regAction {
	return &regAction{
		spinner: view.Spinner(),
		app:     a,
		regData: data,
	}
}

func (rp *regAction) ActionCmd() tea.Msg {
	ctx := context.Background()
	token, err := rp.app.Conn.Register(ctx, rp.regData.Login, rp.regData.Password)
	if err != nil {
		return RegActionFail{Err: errors.New("неверно введены данные")}
	}
	rp.app.Conn.SetToken(token)
	return RegActionSuccess{}
}

func (rp *regAction) Init() tea.Cmd {
	return tea.Batch(rp.spinner.Tick, rp.ActionCmd)
}

func (rp *regAction) Exec() (tea.Model, tea.Cmd) {
	cmd := rp.Init()
	return rp, cmd
}

func (rp *regAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RegActionSuccess:
		return NewLK(rp.app).Exec()

	case RegActionFail:
		return NewReg(rp.app).WithError(msg.Err).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return rp, tea.Quit
		default:
			return rp, nil
		}

	default:
		var cmd tea.Cmd
		rp.spinner, cmd = rp.spinner.Update(msg)
		return rp, cmd
	}
}

func (rp *regAction) View() string {
	str := view.BreakLine().Two()
	str += rp.spinner.View()
	str += " Минутку..." + view.BreakLine().One()
	str += view.Quit()
	return str
}
