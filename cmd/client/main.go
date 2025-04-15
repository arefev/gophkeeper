package main

import (
	"fmt"
	"log"

	"github.com/arefev/gophkeeper/internal/client/app"
	"github.com/arefev/gophkeeper/internal/client/connection"
	"github.com/arefev/gophkeeper/internal/client/tui/step"
	"github.com/arefev/gophkeeper/internal/logger"
)

func main() {
	l, err := logger.Build("debug")
	if err != nil {
		log.Fatal("logger build fail")
	}

	conn := connection.NewGRPCClient(l)
	if err = conn.Connect(":3200"); err != nil {
		log.Fatalf("connect failed: %s", err.Error())
	}

	defer func() {
		if err = conn.Close(); err != nil {
			log.Fatalf("connect close failed: %s", err.Error())
		}
	}()

	a := app.NewApp(conn)
	_, err = step.NewStart(a).NewProgram().Run()
	if err != nil {
		log.Fatalf("app stopped with error: %s", err.Error())
	}

	fmt.Println("app stopped")
}
