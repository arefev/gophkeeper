package step

import (
	"errors"
	"time"

	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginActionSuccess struct {
	AuthToken string
}

type LoginActionFail struct {
	Err error
}

type loginAction struct {
	loginData *model.LoginData
	spinner   spinner.Model
}

func NewLoginAction(data *model.LoginData) *loginAction {
	return &loginAction{
		spinner:   view.Spinner(),
		loginData: data,
	}
}

func (la *loginAction) ActionCmd() tea.Msg {
	time.Sleep(time.Second * 2)
	if la.loginData.Login == "" {
		return LoginActionFail{Err: errors.New("неверный логин/пароль")}
	}

	return LoginActionSuccess{AuthToken: "rkjjfhrehgehrgkhf234231421jeefewf"}
}

func (la *loginAction) Init() tea.Cmd {
	return tea.Batch(la.spinner.Tick, la.ActionCmd)
}

func (la *loginAction) Exec() (tea.Model, tea.Cmd) {
	return la, la.Init()
}

func (la *loginAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoginActionSuccess:
		return NewLogin().Exec()
	case LoginActionFail:
		return NewLogin().WithError(msg.Err).Exec()
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
