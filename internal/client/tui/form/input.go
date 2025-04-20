package form

import (
	"github.com/arefev/gophkeeper/internal/client/tui/style"
	"github.com/charmbracelet/bubbles/textinput"
)

const (
	CharLimit int = 32
)

type Input struct {
	model *textinput.Model
	code  string
}

func NewInput(code, placeholder string, isFocused, isPwd bool) *Input {
	i := &Input{code: code}

	t := textinput.New()
	t.Cursor.Style = style.CursorStyle
	t.CharLimit = CharLimit
	t.Placeholder = placeholder

	if isFocused {
		t.Focus()
		t.PromptStyle = style.FocusedStyle
		t.TextStyle = style.FocusedStyle
	}

	if isPwd {
		t.EchoMode = textinput.EchoPassword
		t.EchoCharacter = '•'
	}

	i.model = &t
	return i
}

func (i *Input) Code() string {
	return i.code
}

func (i *Input) Model() *textinput.Model {
	return i.model
}

func (i *Input) SetModel(model *textinput.Model) {
	i.model = model
}

func (i *Input) SetCharsLimit(limit int) *Input {
	i.model.CharLimit = limit
	return i
}
