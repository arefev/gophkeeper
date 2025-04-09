package step

import (
	"time"

	"github.com/arefev/gophkeeper/internal/client/pipeline/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginProcessingMsg int

type errMsg error

type loginProcessing struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func NewLoginProcessing() loginProcessing {
	return loginProcessing{spinner: view.Spinner()}
}

func (lp loginProcessing) ProcessingCmd() tea.Msg {
	time.Sleep(time.Second * 2)
	return LoginProcessingMsg(1)
}

func (lp loginProcessing) Init() tea.Cmd {
	return tea.Batch(lp.spinner.Tick, lp.ProcessingCmd)
}

func (lp loginProcessing) NewProgram() *tea.Program {
	return tea.NewProgram(lp)
}

func (lp loginProcessing) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case LoginProcessingMsg:
		login := NewLogin().WithError("Неверный логин/пароль")
		return login, login.Init()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			lp.quitting = true
			return lp, tea.Quit
		default:
			return lp, nil
		}

	case errMsg:
		lp.err = msg
		return lp, nil

	default:
		var cmd tea.Cmd
		lp.spinner, cmd = lp.spinner.Update(msg)
		return lp, cmd
	}
}

func (lp loginProcessing) View() string {
	str := view.Break(2)
	str += lp.spinner.View()
	str += " Минутку..." + view.Break(1)
	str += view.Quit()
	return str
}
