package step

import (
	"fmt"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/tui/style"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	return nil
}

func (lt *list) Exec() (tea.Model, tea.Cmd) {
	cmd := lt.Init()
	return lt, cmd
}

func (lt *list) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if msg, ok := msg.(tea.KeyMsg); ok {
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
		{Title: "Uuid", Width: 5},
		{Title: "Тип", Width: 15},
		{Title: "Имя", Width: 20},
		{Title: "Дата создания", Width: 15},
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

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return t
}
