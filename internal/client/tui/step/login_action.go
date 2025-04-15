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

type LoginActionSuccess struct {
}

type LoginActionFail struct {
	Err error
}

type loginAction struct {
	app       *app.App
	loginData *model.LoginData
	spinner   spinner.Model
}

func NewLoginAction(data *model.LoginData, a *app.App) *loginAction {
	return &loginAction{
		spinner:   view.Spinner(),
		loginData: data,
		app:       a,
	}
}

func (la *loginAction) ActionCmd() tea.Msg {
	ctx := context.Background()
	// TODO: validation needed
	token, err := la.app.Conn.Login(ctx, la.loginData.Login, la.loginData.Password)
	if err != nil {
		return LoginActionFail{Err: errors.New("неверный логин/пароль")}
	}
	la.app.Conn.SetToken(token)
	return LoginActionSuccess{}
}

func (la *loginAction) Init() tea.Cmd {
	return tea.Batch(la.spinner.Tick, la.ActionCmd)
}

func (la *loginAction) Exec() (tea.Model, tea.Cmd) {
	cmd := la.Init()
	return la, cmd
}

func (la *loginAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoginActionSuccess:
		return NewLK(la.app).Exec()

	case LoginActionFail:
		return NewLogin(la.app).WithError(msg.Err).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return la, tea.Quit

		default:
			return la, nil
		}

	default:
		var cmd tea.Cmd
		la.spinner, cmd = la.spinner.Update(msg)
		return la, cmd
	}
}

func (la *loginAction) View() string {
	str := view.BreakLine().Two()
	str += la.spinner.View()
	str += " Минутку..." + view.BreakLine().One()
	str += view.Quit()
	return str
}
