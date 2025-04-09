package step_auth

import (
	"time"

	"github.com/arefev/gophkeeper/internal/client/pipeline/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type ProcessingMsg int

type errMsg error

type processing struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func NewProcessing() processing {
	return processing{spinner: view.Spinner()}
}

func (p processing) ProcessingCmd() tea.Msg {
	time.Sleep(time.Second * 2)
	return ProcessingMsg(1)
}

func (p processing) Init() tea.Cmd {
	return tea.Batch(p.spinner.Tick, p.ProcessingCmd)
}

func (p processing) NewProgram() *tea.Program {
	return tea.NewProgram(p)
}

func (p processing) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ProcessingMsg:
		login := NewLogin().WithError("Неверный логин/пароль")
		return login, login.Init()
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			p.quitting = true
			return p, tea.Quit
		default:
			return p, nil
		}
	

	case errMsg:
		p.err = msg
		return p, nil

	default:
		var cmd tea.Cmd
		p.spinner, cmd = p.spinner.Update(msg)
		return p, cmd
	}
}

func (p processing) View() string {
	str := view.Break(2)
	str += p.spinner.View()
	str += " Минутку..." + view.Break(1)
	str += view.Quit()
	return str
}
