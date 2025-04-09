package main

import (
	"os"

	"github.com/arefev/gophkeeper/internal/client/pipeline/step"
)

func main() {
	_, err := step.NewStart().NewProgram().Run()
	if err != nil {
		os.Exit(0)
	}
}
