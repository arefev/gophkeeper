package app

import "context"

type Connection interface {
	Register(ctx context.Context, login, pwd string) error
}

type App struct {
	Conn Connection
}

func NewApp(conn Connection) *App {
	return &App{
		Conn: conn,
	}
}
