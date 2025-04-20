package step

import (
	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	tea "github.com/charmbracelet/bubbletea"
)

type lkTypes struct {
	choice  string
	app     *app.App
	choices []string
	cursor  int
}

func NewLKTypes(a *app.App) *lkTypes {
	return &lkTypes{
		choices: []string{"Логин/пароль", "Банковская карта", "Файл"},
		app:     a,
	}
}

func (lkt *lkTypes) Init() tea.Cmd {
	return lkt.app.Conn.CheckTokenCmd
}

func (lkt *lkTypes) Exec() (tea.Model, tea.Cmd) {
	cmd := lkt.Init()
	return lkt, cmd
}

func (lkt *lkTypes) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case connection.CheckAuthFail:
		return NewStart(lkt.app).Exec()

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return NewLK(lkt.app).Exec()

		case tea.KeyCtrlC:
			return lkt, tea.Quit

		case tea.KeyEnter:
			lkt.choice = lkt.choices[lkt.cursor]

			switch lkt.choices[lkt.cursor] {
			case "Логин/пароль":
				return NewLKFormCreds(lkt.app).Exec()

			case "Банковская карта":
				return NewLKFormBank(lkt.app).Exec()

			default:
				return lkt, tea.Quit
			}

		case tea.KeyDown:
			lkt.cursor++
			if lkt.cursor >= len(lkt.choices) {
				lkt.cursor = 0
			}

		case tea.KeyUp:
			lkt.cursor--
			if lkt.cursor < 0 {
				lkt.cursor = len(lkt.choices) - 1
			}

		default:
			return lkt, nil
		}
	}

	return lkt, nil
}

func (lkt *lkTypes) View() string {
	str := view.Title("Какие данные вы хотите отправить?")

	for i := range lkt.choices {
		if lkt.cursor == i {
			str += "(•) "
		} else {
			str += "( ) "
		}
		str += lkt.choices[i] + view.BreakLine().One()
	}
	str += view.Quit() + view.ToPrevScreen()

	return str
}
