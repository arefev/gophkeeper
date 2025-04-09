package step

import (
	step_auth "github.com/arefev/gophkeeper/internal/client/pipeline/step/auth"
	step_reg "github.com/arefev/gophkeeper/internal/client/pipeline/step/reg"
	"github.com/arefev/gophkeeper/internal/client/pipeline/view"
	tea "github.com/charmbracelet/bubbletea"
)

type Start struct {
	cursor  int
	choice  string
	choices []string
}

func NewStart() Start {
	return Start{
		choices: []string{"Авторизация", "Регистрация"},
	}
}

func (s Start) Init() tea.Cmd {
	return nil
}

func (s Start) NewProgram() *tea.Program {
	return tea.NewProgram(s, tea.WithAltScreen())
}

func (s Start) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return s, tea.Quit

		case tea.KeyEnter:
			s.choice = s.choices[s.cursor]

			switch s.choices[s.cursor] {
			case "Авторизация":
				login := step_auth.NewLogin()
				return login, login.Init()
			case "Регистрация":
				return step_reg.NewReg(), nil
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
		}
	}

	return s, nil
}

func (s Start) View() string {
	str := view.Title("Добро пожаловать ⭐")

	for i := range s.choices {
		if s.cursor == i {
			str += "(•) "
		} else {
			str += "( ) "
		}
		str += s.choices[i] + view.Break(1)
	}
	str += view.Quit()

	return str
}
