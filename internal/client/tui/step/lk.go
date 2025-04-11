package step

import (
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	tea "github.com/charmbracelet/bubbletea"
)

type lk struct {
	choice  string
	choices []string
	cursor  int
}

func NewLK() *lk {
	return &lk{
		choices: []string{"Получить данные", "Загрузить данные"},
	}
}

func (lk *lk) Init() tea.Cmd {
	return nil
}

func (lk *lk) Exec() (tea.Model, tea.Cmd) {
	cmd := lk.Init()
	return lk, cmd
}

func (lk *lk) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyCtrlC:
			return lk, tea.Quit

		case tea.KeyEnter:
			lk.choice = lk.choices[lk.cursor]

			switch lk.choices[lk.cursor] {
			case "Получить данные":
				return NewList().Exec()

			case "Загрузить данные":
				return lk, nil

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
