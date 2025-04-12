package step

import (
	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/tui/view"
	tea "github.com/charmbracelet/bubbletea"
)

type Start struct {
	choice  string
	app     *app.App
	choices []string
	cursor  int
}

func NewStart(a *app.App) *Start {
	return &Start{
		choices: []string{"Авторизация", "Регистрация"},
		app:     a,
	}
}

func (s *Start) Init() tea.Cmd {
	return nil
}

func (s *Start) Exec() (tea.Model, tea.Cmd) {
	cmd := s.Init()
	return s, cmd
}

func (s *Start) NewProgram() *tea.Program {
	return tea.NewProgram(s, tea.WithAltScreen())
}

func (s *Start) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.Type {
		case tea.KeyCtrlC:
			return s, tea.Quit

		case tea.KeyEnter:
			s.choice = s.choices[s.cursor]

			switch s.choices[s.cursor] {
			case "Авторизация":
				return NewLogin(s.app).Exec()

			case "Регистрация":
				return NewReg(s.app).Exec()

			default:
				return s, tea.Quit
			}

		case tea.KeyDown:
			s.cursor++
			if s.cursor >= len(s.choices) {
				s.cursor = 0
			}

		case tea.KeyUp:
			s.cursor--
			if s.cursor < 0 {
				s.cursor = len(s.choices) - 1
			}

		default:
			return s, nil
		}
	}

	return s, nil
}

func (s *Start) View() string {
	str := view.Title("Добро пожаловать ⭐")

	for i := range s.choices {
		if s.cursor == i {
			str += "(•) "
		} else {
			str += "( ) "
		}
		str += s.choices[i] + view.BreakLine().One()
	}
	str += view.Quit()

	return str
}
