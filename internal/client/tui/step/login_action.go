package step

import (
	"errors"
	"time"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginActionSuccess struct {
	AuthData model.AuthData
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
	const s = 2
	time.Sleep(time.Second * s)
	if la.loginData.Login == "" || la.loginData.Password == "" {
		return LoginActionFail{Err: errors.New("неверный логин/пароль")}
	}

	d := model.AuthData{
		Token: "rkjjfhrehgehrgkhf234231421jeefewf",
	}

	return LoginActionSuccess{AuthData: d}
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
		return NewLK().Exec()

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
