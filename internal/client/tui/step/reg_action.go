package step

import (
	"errors"
	"time"

	"github.com/arefev/gophkeeper/internal/client/app"
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
	spinner spinner.Model
	app     *app.App
}

func NewRegAction(app *app.App) *regAction {
	return &regAction{
		spinner: view.Spinner(),
		app:     app,
	}
}

func (rp *regAction) ProcessingCmd() tea.Msg {
	const s = 2
	time.Sleep(time.Second * s)
	return RegActionFail{Err: errors.New("неверно введены данные")}
}

func (rp *regAction) Init() tea.Cmd {
	return tea.Batch(rp.spinner.Tick, rp.ProcessingCmd)
}

func (rp *regAction) Exec() (tea.Model, tea.Cmd) {
	cmd := rp.Init()
	return rp, cmd
}

func (rp *regAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RegActionSuccess:
		return NewReg(rp.app).Exec()

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
