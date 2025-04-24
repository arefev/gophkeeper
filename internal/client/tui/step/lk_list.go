package step

import (
	"context"
	"fmt"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/style"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type lkList struct {
	app   *app.App
	table table.Model
}

func NewLKList(a *app.App) *lkList {
	return &lkList{
		table: getTable(),
		app:   a,
	}
}

func (lkl *lkList) ActionCmd() tea.Msg {
	ctx := context.Background()
	// TODO: validation needed
	err := lkl.app.Conn.GetList(ctx)
	if err != nil {
		lkl.app.Log.Error("GetList failed", zap.Error(err))
		return nil
	}

	return nil
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

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewLK(lkl.app).Exec()

		case tea.KeyCtrlC:
			return lkl, tea.Quit

		case tea.KeyEnter:
			return lkl, tea.Batch(
				tea.Printf("Let's go to %s!", lkl.table.SelectedRow()[1]),
			)

		default:
			lkl.table, cmd = lkl.table.Update(msg)
			return lkl, cmd
		}
	}

	lkl.table, cmd = lkl.table.Update(msg)
	return lkl, cmd
}

func (lkl *lkList) View() string {
	str := view.Title("Выберите данные 📃")
	str += style.BorderStyle.Render(lkl.table.View()) + view.BreakLine().One()
	str += fmt.Sprintf("Всего строк: %d", len(lkl.table.Rows()))
	str += view.BreakLine().One()
	str += view.Quit() + view.ToPrevScreen()
	return str
}

func getTable() table.Model {
	columns := []table.Column{
		{Title: "Uuid", Width: style.ColumnWidthS},
		{Title: "Тип", Width: style.ColumnWidthM},
		{Title: "Имя", Width: style.ColumnWitdthL},
		{Title: "Дата создания", Width: style.ColumnWidthM},
	}

	rows := []table.Row{
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72977", "Карта", "Сбер", "01.03.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72988", "Файл", "load.exe", "02.02.2025"},
		{"1e5b491b-39a3-40d2-92ed-adabfbb72999", "Логин/пароль", "https://habr.com", "03.03.2025"},
	}

	return view.Table(columns, rows)
}
