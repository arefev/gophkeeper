package step

import (
	"context"
	"errors"
	"fmt"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/model"
	"github.com/arefev/gophkeeper/internal/client/tui/style"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type lkList struct {
	app     *app.App
	table   table.Model
	list    *[]model.MetaListData
	spinner spinner.Model
}

type IsListLoaded bool
type ListActionFail struct {
	Err error
}

func NewLKList(a *app.App) *lkList {
	return &lkList{
		app:     a,
		spinner: view.Spinner(),
	}
}

func (lkl *lkList) ActionCmd() tea.Msg {
	ctx := context.Background()
	list, err := lkl.app.Conn.GetList(ctx)
	if err != nil {
		lkl.app.Log.Error("GetList failed", zap.Error(err))
		return ListActionFail{Err: errors.New("не удалось загрузить данные")}
	}

	lkl.list = list
	lkl.table = lkl.getTable()
	return IsListLoaded(true)
}

func (lkl *lkList) Init() tea.Cmd {
	return tea.Batch(lkl.app.Conn.CheckTokenCmd, lkl.ActionCmd)
}

func (lkl *lkList) Exec() (tea.Model, tea.Cmd) {
	cmd := lkl.Init()
	return lkl, cmd
}

func (lkl *lkList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(lkl.app).Exec()

	case ListActionFail:
		return NewLK(lkl.app).WithError(msg.Err).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewLK(lkl.app).Exec()

		case tea.KeyCtrlC:
			return lkl, tea.Quit

		case tea.KeyEnter:
			return NewDownloadAction(lkl.table.SelectedRow()[0], lkl.app).Exec()

		default:
			lkl.table, cmd = lkl.table.Update(msg)
			return lkl, cmd
		}
	}

	lkl.table, cmd = lkl.table.Update(msg)
	return lkl, cmd
}

func (lkl *lkList) View() string {
	if lkl.list == nil {
		str := view.BreakLine().Two()
		str += lkl.spinner.View()
		str += " Минутку..." + view.BreakLine().One()
		str += view.Quit()
		return str
	}

	str := view.Title("Выберите данные 📃")
	str += style.BorderStyle.Render(lkl.table.View()) + view.BreakLine().One()
	str += fmt.Sprintf("Всего строк: %d", len(lkl.table.Rows()))
	str += view.BreakLine().One()
	str += view.Quit() + view.ToPrevScreen()
	return str
}

func (lkl *lkList) getTable() table.Model {
	columns := []table.Column{
		{Title: "Uuid", Width: style.ColumnWidthS},
		{Title: "Тип", Width: style.ColumnWidthM},
		{Title: "Имя", Width: style.ColumnWidthL},
		{Title: "Файл", Width: style.ColumnWidthL},
		{Title: "Дата создания", Width: style.ColumnWidthL},
	}

	rows := make([]table.Row, 0, len(*lkl.list))
	for _, item := range *lkl.list {
		row := table.Row{
			item.UUID,
			item.Type,
			item.Name,
			item.File,
			item.Date,
		}

		rows = append(rows, row)
	}

	return view.Table(columns, rows)
}
