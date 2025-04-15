package step

import (
	"fmt"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/style"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
)

type list struct {
	app   *app.App
	table table.Model
}

func NewList(a *app.App) *list {
	return &list{
		table: getTable(),
		app:   a,
	}
}

func (lt *list) Init() tea.Cmd {
	return lt.app.Conn.CheckTokenCmd
}

func (lt *list) Exec() (tea.Model, tea.Cmd) {
	cmd := lt.Init()
	return lt, cmd
}

func (lt *list) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(lt.app).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewLK(lt.app).Exec()

		case tea.KeyCtrlC:
			return lt, tea.Quit

		case tea.KeyEnter:
			return lt, tea.Batch(
				tea.Printf("Let's go to %s!", lt.table.SelectedRow()[1]),
			)

		default:
			lt.table, cmd = lt.table.Update(msg)
			return lt, cmd
		}
	}

	lt.table, cmd = lt.table.Update(msg)
	return lt, cmd
}

func (lt *list) View() string {
	str := view.Title("Выберите данные 📃")
	str += style.BorderStyle.Render(lt.table.View()) + view.BreakLine().One()
	str += fmt.Sprintf("Всего строк: %d", len(lt.table.Rows()))
	str += view.BreakLine().One()
	str += view.Quit() + view.ToStart()
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
