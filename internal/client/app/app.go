package app

import (
	"context"

	"github.com/arefev/gophkeeper/internal/client/tui/model"
	tea "github.com/charmbracelet/bubbletea"
	"go.uber.org/zap"
)

type Connection interface {
	Register(ctx context.Context, login, pwd string) (string, error)
	Login(ctx context.Context, login, pwd string) (string, error)
	FileUpload(ctx context.Context, path, metaName, metaType string) error
	TextUpload(ctx context.Context, txt []byte, metaName, metaType string) error
	GetList(ctx context.Context) (*[]model.MetaListData, error)
	FileDownload(ctx context.Context, uuid string) (string, error)
	Delete(ctx context.Context, uuid string) error
	SetToken(t string)
	CheckTokenCmd() tea.Msg
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
