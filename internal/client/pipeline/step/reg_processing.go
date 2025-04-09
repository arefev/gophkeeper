package step

import (
	"time"

	"github.com/arefev/gophkeeper/internal/client/pipeline/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type RegProcessingMsg int

type regProcessing struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func NewRegProcessing() regProcessing {
	return regProcessing{spinner: view.Spinner()}
}

func (rp regProcessing) ProcessingCmd() tea.Msg {
	time.Sleep(time.Second * 2)
	return RegProcessingMsg(1)
}

func (rp regProcessing) Init() tea.Cmd {
	return tea.Batch(rp.spinner.Tick, rp.ProcessingCmd)
}

func (rp regProcessing) Exec() (tea.Model, tea.Cmd) {
	return rp, rp.Init()
}

func (rp regProcessing) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case RegProcessingMsg:
		login := NewReg().WithError("Неверно заполнены данные")
		return login, login.Init()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			rp.quitting = true
			return rp, tea.Quit
		default:
			return rp, nil
		}

	case errMsg:
		rp.err = msg
		return rp, nil

	default:
		var cmd tea.Cmd
		rp.spinner, cmd = rp.spinner.Update(msg)
		return rp, cmd
	}
}

func (rp regProcessing) View() string {
	str := view.Break(2)
	str += rp.spinner.View()
	str += " Минутку..." + view.Break(1)
	str += view.Quit()
	return str
}
