package interceptor

import "github.com/arefev/gophkeeper/internal/server/application"

type middleware struct {
	app *application.App
}

func NewMiddleware(app *application.App) *middleware {
	return &middleware{
		app: app,
	}
}
