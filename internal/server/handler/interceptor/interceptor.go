package interceptor

import "github.com/arefev/gophkeeper/internal/server/application"

type interceptor struct {
	app *application.App
}

func New(app *application.App) *interceptor {
	return &interceptor{
		app: app,
	}
}
