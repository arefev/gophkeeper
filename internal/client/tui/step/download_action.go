package step

import (
	"context"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type DownloadActionSuccess struct {
	FilePath string
}

type DownloadActionFail struct {
	Err error
}

type downloadAction struct {
	app     *app.App
	uuid    string
	spinner spinner.Model
}

func NewDownloadAction(uuid string, a *app.App) *downloadAction {
	return &downloadAction{
		spinner: view.Spinner(),
		uuid:    uuid,
		app:     a,
	}
}

func (da *downloadAction) ActionCmd() tea.Msg {
	path, err := da.app.Conn.FileDownload(context.Background(), da.uuid)
	if err != nil {
		return DownloadActionFail{Err: err}
	}
	return DownloadActionSuccess{FilePath: path}
}

func (da *downloadAction) Init() tea.Cmd {
	return tea.Batch(da.spinner.Tick, da.app.Conn.CheckTokenCmd, da.ActionCmd)
}

func (da *downloadAction) Exec() (tea.Model, tea.Cmd) {
	cmd := da.Init()
	return da, cmd
}

func (da *downloadAction) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(da.app).Exec()

	case DownloadActionSuccess:
		return NewLKList(da.app).WithMsg("Файл успешно сохранён - " + msg.FilePath).Exec()

	case DownloadActionFail:
		return NewLKList(da.app).WithError(msg.Err).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return da, tea.Quit

		default:
			return da, nil
		}

	default:
		var cmd tea.Cmd
		da.spinner, cmd = da.spinner.Update(msg)
		return da, cmd
	}
}

func (da *downloadAction) View() string {
	str := view.BreakLine().Two()
	str += da.spinner.View()
	str += " Минутку..." + view.BreakLine().One()
	str += view.Quit()
	return str
}
