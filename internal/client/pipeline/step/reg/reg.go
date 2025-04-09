package step_reg


import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)


type reg struct {
	cursor  int
	choice  string
	choices []string
}

func NewReg() reg {
	return reg{
		choices: []string{"Шаг reg 1", "Шаг reg 2"},
	}
}

func (r reg) NewProgram() *tea.Program {
	return tea.NewProgram(r)
}

func (r reg) Init() tea.Cmd {
	return nil
}

func (r reg) Choice() string {
	return r.choice
}

func (r reg) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return r, tea.Quit

		case "enter":
			// Send the choice on the channel and exit.
			r.choice = r.choices[r.cursor]
			return r, tea.Quit

		case "down", "j":
			r.cursor++
			if r.cursor >= len(r.choices) {
				r.cursor = 0
			}

		case "up", "k":
			r.cursor--
			if r.cursor < 0 {
				r.cursor = len(r.choices) - 1
			}
		}
	}

	return r, nil
}

func (r reg) View() string {
	s := strings.Builder{}
	s.WriteString("Авторизуйтесь\n\n")

	for i := 0; i < len(r.choices); i++ {
		if r.cursor == i {
			s.WriteString("(•) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(r.choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\n(press q to quit)\n")

	return s.String()
}
