package main

import (
	"os"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/tui/step"
)

func main() {
	a := app.NewApp()
	_, err := step.NewStart(a).NewProgram().Run()
	if err != nil {
		os.Exit(0)
	}
}
