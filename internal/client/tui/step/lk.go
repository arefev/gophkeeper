package step

import (
	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	tea "github.com/charmbracelet/bubbletea"
)

type lk struct {
	choice  string
	app     *app.App
	choices []string
	cursor  int
}

func NewLK(a *app.App) *lk {
	return &lk{
		choices: []string{"Получить данные", "Загрузить данные"},
		app:     a,
	}
}

func (lk *lk) Init() tea.Cmd {
	return lk.app.Conn.CheckTokenCmd
}

func (lk *lk) Exec() (tea.Model, tea.Cmd) {
	cmd := lk.Init()
	return lk, cmd
}

func (lk *lk) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(lk.app).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return lk, tea.Quit

		case tea.KeyEnter:
			lk.choice = lk.choices[lk.cursor]

			switch lk.choices[lk.cursor] {
			case "Получить данные":
				return NewList(lk.app).Exec()

			case "Загрузить данные":
				return NewLKTypes(lk.app).Exec()

			default:
				return lk, tea.Quit
			}

		case tea.KeyDown:
			lk.cursor++
			if lk.cursor >= len(lk.choices) {
				lk.cursor = 0
			}

		case tea.KeyUp:
			lk.cursor--
			if lk.cursor < 0 {
				lk.cursor = len(lk.choices) - 1
			}

		default:
			return lk, nil
		}
	}

	return lk, nil
}

func (lk *lk) View() string {
	str := view.Title("Личный кабинет 🔑")

	for i := range lk.choices {
		if lk.cursor == i {
			str += "(•) "
		} else {
			str += "( ) "
		}
		str += lk.choices[i] + view.BreakLine().One()
	}
	str += view.Quit()

	return str
}
