package app

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type Connection interface {
	Register(ctx context.Context, login, pwd string) (string, error)
	Login(ctx context.Context, login, pwd string) (string, error)
	SetToken(t string)
	CheckTokenCmd() tea.Msg
	FileUpload(ctx context.Context, creds []byte) error
}

type App struct {
	Conn Connection
	Log  *zap.Logger
}

func NewApp(conn Connection, log *zap.Logger) *App {
	return &App{
		Conn: conn,
		Log:  log,
	}
}
