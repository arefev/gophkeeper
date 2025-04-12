package main

import (
	"fmt"
	"os"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/tui/step"
)

func main() {
	a := app.NewApp()
	_, err := step.NewStart(a).NewProgram().Run()
	if err != nil {
		fmt.Println("app stopped with error: ", err.Error())
		os.Exit(0)
	}

	fmt.Println("app stopped")
}
